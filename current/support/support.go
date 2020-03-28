package support

import (
	"fmt"
	"os"
	"time"
)

func ErrorIt(errormsg string, errornum int) {
	currentTime := time.Now()
	fmt.Println(currentTime.Format("2006-01-02 15:04:05"), ": error code:", errornum, "; Error: ", errormsg)
	os.Exit(errornum)
}

func DebugIt(debugmsg string) {
	currentTime := time.Now()
	//isDebug := os.Getenv("OTTODEBUG")
	//fmt.Println()
	/*for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}*/
	if os.Getenv("OTTODEBUG") == "1" {
		//if os.Getenv("OttoDebug") != "" {
		var dbgmsg string
		dbgmsg = ": Debug: " + debugmsg
		fmt.Println(currentTime.Format("2006-01-02 15:04:05"), dbgmsg)
	}
}

func NormalPrint(normalprt string) {
	currentTime := time.Now()
	var msg string
	msg = ": " + normalprt
	fmt.Println(currentTime.Format("2006-01-02 15:04:05"), msg)
}
