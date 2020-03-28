package ottotf

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"hangar.bpfx.org/otto/otto2/support"

	pongo "github.com/flosch/pongo2"
)

//type Comp2 map[string]string

/*
func WriteTF(labelname string,
	vspherename string,
	servername string,
	templatename string,
	networkname string,
	ipaddress string,
	gatewayaddress string,
	bitmask string,
	nmgorg string,
	nmgowner string,
	nmgapplication string,
	templatefile string,
	vcluster string,
	dcname string,
	storagecluster string,
	outfile string) (err2 error) {
	templ := pongo.Must(pongo.FromFile(templatefile))
	w, err1 := os.Create(outfile)
	if err1 != nil {
		support.DebugIt("Error in WriteTF")
	}
	err := templ.ExecuteWriter(pongo.Context{"labelname": labelname, "vspherename": vspherename, "servername": servername, "templatename": templatename, "networkname": networkname, "ipaddress": ipaddress, "gatewayaddress": gatewayaddress, "bitmask": bitmask, "nmgorg": nmgorg, "nmgowner": nmgowner, "nmgapplication": nmgapplication}, w)
	return (err)
}
*/

func WriteTF(server map[string]string, outfile string, templatefile string) (err2 error) {
	support.DebugIt("Template: " + templatefile)
	var templ = pongo.Must(pongo.FromFile(templatefile))
	w, err1 := os.Create(outfile)
	if err1 != nil {
		support.DebugIt("Error in WriteTF")
	}
	//support.DebugIt(server)
	err := templ.ExecuteWriter(pongo.Context{"labelname": server["labelname"],
		"vspherename":    server["vspherename"],
		"servername":     server["servername"],
		"templatename":   server["templatename"],
		"vcluster":       server["vcluster"],
		"dcname":         server["dcname"],
		"networkname":    server["networkname"],
		"storagecluster": server["storagecluster"],
		"ipaddress":      server["ipaddress"],
		"gatewayaddress": server["gatewayaddress"],
		"bitmask":        server["bitmask"],
		"nmgorg":         server["nmgorg"],
		"nmgowner":       server["nmgowner"],
		"nmgapplication": server["nmgapplication"]}, w)
	return (err)
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func CopyFile(source string, dest string) {
	destfile, err := os.Create(dest)
	defer destfile.Close()
	sourcefile, err := os.Open(source)
	defer sourcefile.Close()
	io.Copy(destfile, sourcefile)
	if err != nil {
		fmt.Println("Error")
	}
}

func RunTF(tffile string, servername string, user string, pass string, tempdir string) {
	/*
	 * First, setup the temp directory
	 */
	/*b := String(8)
	var tempdir string
	tempdir = "/tmp/" + string(b)
	os.Mkdir(tempdir, 0755)
	os.Chdir(tempdir)

	CopyFile(tffile, tempdir+"/"+servername)
	*/
	os.Chdir(tempdir)
	cmd := exec.Command("terraform", "init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	CopyFile("/Users/nmrb126/.terraform.tfvars", tempdir+"/terraform.tfvars")

	os.Setenv("TF_VAR_user=", user)
	os.Setenv("TF_VAR_pass=", pass)
	cmd2 := exec.Command("terraform", "plan")
	out2, err2 := cmd2.CombinedOutput()
	if err2 != nil {
		fmt.Printf("combined out:\n%s\n", string(out2))
		log.Fatalf("cmd.Run() failed with %s\n", err2)
	}
	//fmt.Printf("combined out:\n%s\n", string(out2))

	//	Cleanup(tempdir)
}

func Cleanup(thedir string) {
	os.RemoveAll(thedir)
}
