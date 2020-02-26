//usr/bin/env go run "$0" "$@"; exit "$?"

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"ottonet"
	_ "ottonet"
	"ottotf"
	"support"
)

var debugenabled bool

/*
 * Error 1: missing filename for the xml file
 * Error 2: error opening file
 */

// MxfileSt Main structure of the xml
type MxfileSt struct {
	XMLName xml.Name  `xml:"mxfile"`
	Diagram DiagramSt `xml:"diagram"`
}

// DiagramSt next structure in xml
type DiagramSt struct {
	XMLName      xml.Name       `xml:"diagram"`
	MXGraphModel mxGraphModelSt `xml:"mxGraphModel"`
}

// mxGraphModeSt next structure in xml
type mxGraphModelSt struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Root    RootSt   `xml:"root"`
}

// RootSt next structure in xml
type RootSt struct {
	XMLName        xml.Name       `xml:"root"`
	ObjectList     []ObjectSt     `xml:"object"`
	ConnectionList []ConnectionSt `xml:"mxCell"`
}

// ObjectSt This is what contains what you want, unless you want the connection info which is in another castle.
type ObjectSt struct {
	XMLName        xml.Name `xml:"object"`
	Label          string   `xml:"label,attr"`
	Name           string   `xml:"nmgname,attr"`
	Type           string   `xml:"nmgtype,attr"`
	OSType         string   `xml:"nmgostype,attr"`
	ID             string   `xml:"id,attr"`
	NmgOrg         string   `xml:"nmgorg,attr"`
	NmgOwner       string   `xml:"nmgowner,attr"`
	NmgApplication string   `xml:"nmgapplication,attr"`
}

// ConnectionSt struct containing the connections
type ConnectionSt struct {
	XMLName xml.Name `xml:"mxCell"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}

type Comp map[string]string

/*func Getnetworkdeets(networkname string) (netname string) {
	return (networkname)
}*/

func main() {
	/*
	 * Work on command line stuff first
	 */
	// Todo: send the debug flag to debugIt
	debug := flag.Bool("d", false, "Adds debugging")

	var xmlfilename string
	flag.StringVar(&xmlfilename, "f", "", "XML File")

	var templatefilename string
	flag.StringVar(&templatefilename, "t", "", "VM Template File")

	var esxuser string
	flag.StringVar(&esxuser, "user", "", "Username for ESX")

	var esxpass string
	flag.StringVar(&esxpass, "password", "", "Password for ESX")

	flag.Parse()
	fmt.Println("Debug:", *debug)

	debugenabled = *debug

	if esxuser == "" {
		fmt.Print("Enter Username: ")
		byteUser, err := terminal.ReadPassword(0)
		if err == nil {
			fmt.Println("\nError User")
		}
		esxuser = string(byteUser)
	}

	if esxpass == "" {
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(0)
		if err == nil {
			fmt.Println("\nPassword typed: " + string(bytePassword))
		}
		esxpass = string(bytePassword)
	}

	// Todo: Convert to map
	type RouterSVISt struct {
		Label       string
		NetworkName string
		ID          string
	}
	// Todo: Convert to map
	type ServerSt struct {
		Label          string
		ServerName     string
		IsOnNetwork    string
		OSType         string
		ID             string
		NmgOrg         string
		NmgOwner       string
		NmgApplication string
	}
	// Todo: Convert to map
	type ConnectSt struct {
		ID         string
		Source     string
		Target     string
		TargetName string
	}

	var ServerList []ServerSt
	var RouterSVIList []RouterSVISt
	var ConnectList []ConnectSt

	xmlFile, err := os.Open(xmlfilename)

	if err != nil {
		support.ErrorIt(err.Error(), 2)
	} else {
		support.DebugIt("Opened " + xmlfilename)
	}

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var mxfile MxfileSt

	xml.Unmarshal(byteValue, &mxfile)

	/*
	 * Build lists of the servers, routers (actually which network)
	 */
	for i := 0; i < len(mxfile.Diagram.MXGraphModel.Root.ObjectList); i++ {
		/*
		 * List of the servers
		 */
		if mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Type == "server" {
			tmp := ServerSt{ServerName: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name,
				OSType:         mxfile.Diagram.MXGraphModel.Root.ObjectList[i].OSType,
				ID:             mxfile.Diagram.MXGraphModel.Root.ObjectList[i].ID,
				Label:          mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Label,
				NmgOrg:         mxfile.Diagram.MXGraphModel.Root.ObjectList[i].NmgOrg,
				NmgApplication: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].NmgApplication,
				NmgOwner:       mxfile.Diagram.MXGraphModel.Root.ObjectList[i].NmgOwner}
			support.DebugIt("Adding Entry to ServerList: " + mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name)
			ServerList = append(ServerList, tmp)
		}
		/*
		 * List of RouterSVI (gateways on the routers)
		 */
		if mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Type == "routersvi" {
			tmp2 := RouterSVISt{NetworkName: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name,
				ID: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].ID}
			support.DebugIt("Adding Entry to Router SVI List: " + mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name)
			RouterSVIList = append(RouterSVIList, tmp2)
		}
	}

	/*
	 * List of connections
	 */
	for i := 0; i < len(mxfile.Diagram.MXGraphModel.Root.ConnectionList); i++ {
		srctmp := mxfile.Diagram.MXGraphModel.Root.ConnectionList[i].Source
		tgttmp := mxfile.Diagram.MXGraphModel.Root.ConnectionList[i].Target
		var tgtname string
		if len(srctmp) > 0 && len(tgttmp) > 0 {
			for j := 0; j < len(RouterSVIList); j++ {
				if tgttmp == RouterSVIList[j].ID {
					tgtname = RouterSVIList[j].NetworkName
				}
			}
			tmp3 := ConnectSt{ID: mxfile.Diagram.MXGraphModel.Root.ConnectionList[i].ID,
				Source:     srctmp,
				Target:     tgttmp,
				TargetName: tgtname}
			support.DebugIt("Adding Entry to Connect List: " + tmp3.Source + " -> " + tmp3.Target)
			ConnectList = append(ConnectList, tmp3)
		}
	}

	ServerListNum := len(ServerList)
	RouterSVIListNum := len(RouterSVIList)
	ConnectListNum := len(ConnectList)

	//var SLmsg string
	SLmsg := fmt.Sprintf("Found %d Servers", ServerListNum)
	support.NormalPrint(SLmsg)

	//var RSLmsg string
	RSLmsg := fmt.Sprintf("Found %d Router SVI's", RouterSVIListNum)
	support.NormalPrint(RSLmsg)

	CLmsg := fmt.Sprintf("Found %d Connections", ConnectListNum)
	support.NormalPrint(CLmsg)

	for i := 0; i < len(ServerList); i++ {
		support.DebugIt("Server: " + ServerList[i].ServerName)
	}

	for i := 0; i < len(RouterSVIList); i++ {
		support.DebugIt("Router SVI: " + RouterSVIList[i].NetworkName)
	}

	/*
	 * Logic for this:
	 * Look through the connections, source connects to target
	 * Source is the server
	 * Target is the gateway (network)
	 * Walk through each connection, look at what the source and target
	 * Update the ServerST.IsOnNetwork to reflect which network the device is on
	 */

	for i := 0; i < len(ConnectList); i++ {
		support.DebugIt("Looking at ConnectList" + ConnectList[i].Source)
		for j := 0; j < len(ServerList); j++ {
			support.DebugIt("Looking at ServerList" + ServerList[j].ID)
			if ServerList[j].ID == ConnectList[i].Source {
				ServerList[j].IsOnNetwork = ConnectList[i].TargetName
			}
		}
	}

	if debugenabled {
		for i := 0; i < len(ServerList); i++ {
			fmt.Printf("%+v\n", ServerList[i])
		}
	}

	/*
			 * Now build the tf's
			 * Information we need:
			 * labelname
			 * vspherename
			 * servername
			 * templatename
			 * networkname
			 * ipaddress
			 * gatewayaddress
			 * bitmask
			 * nmgorg
			 * nmgowner
			 * nmgapplication
		   * templatefilename
			 * outputfilename
	*/
	/*
		ottotf.WriteTF(labelname,
			vspherename,
			dcname,
			servername,
			templatename,
			networkname,
			ipaddress,
			gatewayaddress,
			bitmask,
			nmgorg,
			nmgowner,
			nmgapplication,
			templatefilename,
			outputfilename)
	*/

	var ComputeSlice []Comp
	for i := 0; i < len(ServerList); i++ {
		/*
		 * vspherename, clustername will out of Getnetworkdeets
		 */
		var vcenter string
		var dcname string
		var vcluster string
		var storecluster string
		var templatename string

		networkname := ottonet.Getnetworkdeets(ServerList[i].IsOnNetwork)
		fmt.Println(networkname.Name)
		if networkname.Location == "Corp" {
			vcenter = "PLPVCCORP01.corpnmg.net"
			dcname = "PDC Corp"
			if ServerList[i].OSType == "rhel" {
				vcluster = "RHEL Corp"
				storecluster = "VMAX-RHEL-Corp"
				templatename = "rhel7template"
			} else {
				vcluster = "Corp"
				storecluster = "VMAX-Corp"
				templatename = "Windows 2016 Standard_Gold"
			}
		} else {
			vcenter = "PLPLVCCDE01.corpnmg.net"
			dcname = "PDC SDE"
			if ServerList[i].OSType == "rhel" {
				vcluster = "RHEL SDE"
				storecluster = "VMAX-RHEL-SDE"
				templatename = "rhel7template"
			} else {
				vcluster = "SDE"
				storecluster = "VMAX-SDE"
				templatename = "Windows 2016 Standard_Gold"
			}
		}
		var outfile string
		outfile = "/tmp/" + ServerList[i].ServerName + ".tf"
		var ipaddy string
		ipaddy = fmt.Sprintf("172.18.37.20%d", i)
		fmt.Println(ipaddy)
		/*
		 * ipaddres, gatewayaddress and bitmask will need to come from somewhere
		 */
		node1 := Comp{"labelname": ServerList[i].Label,
			"servername":     ServerList[i].ServerName,
			"networkname":    networkname.Name,
			"gatewayaddress": networkname.Gateway,
			"bitmask":        networkname.Bitmask,
			"vspherename":    vcenter,
			"dcname":         dcname,
			"vcluster":       vcluster,
			"storagecluster": storecluster,
			"templatename":   templatename,
			"nmgorg":         ServerList[i].NmgOrg,
			"nmgowner":       ServerList[i].NmgOwner,
			"nmgapplication": ServerList[i].NmgApplication,
			"ipaddress":      ipaddy}
		fmt.Println(node1)
		ComputeSlice = append(ComputeSlice, node1)
		err := ottotf.WriteTF(node1, outfile, templatefilename)
		fmt.Println(err)
		ottotf.RunTF(outfile, ServerList[i].ServerName+".tf", esxuser, esxpass)
	}
}
