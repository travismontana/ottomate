//usr/bin/env go run "$0" "$@"; exit "$?"

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var debugenabled bool = false

/*
 * Error 1: missing filename for the xml file
 * Error 2: error opening file
 */
func errorIt(errormsg string, errornum int) {
	currentTime := time.Now()
	fmt.Println(currentTime.Format("2006-01-02 15:04:05"), ": error code:", errornum, "; Error: ", errormsg)
	os.Exit(errornum)
}

func debugIt(debugmsg string) {
	currentTime := time.Now()
	if debugenabled {
		var dbgmsg string
		dbgmsg = ": Debug: " + debugmsg
		fmt.Println(currentTime.Format("2006-01-02 15:04:05"), dbgmsg)
	}
}

func normalPrint(normalprt string) {
	currentTime := time.Now()
	var msg string
	msg = ": " + normalprt
	fmt.Println(currentTime.Format("2006-01-02 15:04:05"), msg)
}

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
	XMLName xml.Name `xml:"object"`
	Label   string   `xml:"label,attr"`
	Name    string   `xml:"nmgname,attr"`
	Type    string   `xml:"nmgtype,attr"`
	OSType  string   `xml:"nmgostype,attr"`
	ID      string   `xml:"id,attr"`
}

// ConnectionSt struct containing the connections
type ConnectionSt struct {
	XMLName xml.Name `xml:"mxCell"`
	ID      string   `xml:"id,attr"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}

func main() {
	var xmlfilename string

	type RouterSVISt struct {
		Label       string
		NetworkName string
		ID          string
	}

	type ServerSt struct {
		Label       string
		ServerName  string
		IsOnNetwork string
		OSType      string
		ID          string
	}

	type ConnectSt struct {
		ID         string
		Source     string
		Target     string
		TargetName string
	}

	var ServerList []ServerSt
	var RouterSVIList []RouterSVISt
	var ConnectList []ConnectSt

	if len(os.Args) < 2 {
		errorIt("Missing Filename", 1)
	} else {
		xmlfilename = os.Args[1]
		debugIt("Filename: " + xmlfilename)
	}

	xmlFile, err := os.Open(xmlfilename)

	if err != nil {
		errorIt(err.Error(), 2)
	} else {
		debugIt("Opened " + xmlfilename)
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
				OSType: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].OSType,
				ID:     mxfile.Diagram.MXGraphModel.Root.ObjectList[i].ID}
			debugIt("Adding Entry to ServerList: " + mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name)
			ServerList = append(ServerList, tmp)
		}
		/*
		 * List of RouterSVI (gateways on the routers)
		 */
		if mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Type == "routersvi" {
			tmp2 := RouterSVISt{NetworkName: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name,
				ID: mxfile.Diagram.MXGraphModel.Root.ObjectList[i].ID}
			debugIt("Adding Entry to Router SVI List: " + mxfile.Diagram.MXGraphModel.Root.ObjectList[i].Name)
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
			debugIt("Adding Entry to Connect List: " + tmp3.Source + " -> " + tmp3.Target)
			ConnectList = append(ConnectList, tmp3)
		}
	}

	ServerListNum := len(ServerList)
	RouterSVIListNum := len(RouterSVIList)
	ConnectListNum := len(ConnectList)

	//var SLmsg string
	SLmsg := fmt.Sprintf("Found %d Servers", ServerListNum)
	normalPrint(SLmsg)

	//var RSLmsg string
	RSLmsg := fmt.Sprintf("Found %d Router SVI's", RouterSVIListNum)
	normalPrint(RSLmsg)

	CLmsg := fmt.Sprintf("Found %d Connections", ConnectListNum)
	normalPrint(CLmsg)

	for i := 0; i < len(ServerList); i++ {
		debugIt("Server: " + ServerList[i].ServerName)
	}

	for i := 0; i < len(RouterSVIList); i++ {
		debugIt("Router SVI: " + RouterSVIList[i].NetworkName)
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
		debugIt("Looking at ConnectList" + ConnectList[i].Source)
		for j := 0; j < len(ServerList); j++ {
			debugIt("Looking at ServerList" + ServerList[j].ID)
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
}
