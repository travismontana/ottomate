package ottonet

import (
	"fmt"
)
func Getnetworkdeets(networkname string) (netname NetworkSt) {
	var whattoreturn NetworkSt

	tmp := NetworkSt{
	Name: "corpentapps",
	VLANID: "137",
	Gateway: "172.18.37.1",
	Netmask: "255.255.255.0",
	Bitmask: "24",
	Location: "Corp"}
	PDCNetworkList = append(PDCNetworkList,tmp)

	tmp2 := NetworkSt{
	Name: "sdeentapps",
	VLANID: "1137",
	Gateway: "172.30.37.1",
	Netmask: "255.255.255.0",
	Bitmask: "24",
	Location: "SDE"}
	PDCNetworkList = append(PDCNetworkList,tmp2)

	fmt.Println(len(PDCNetworkList))

	for i := 0; i < len(PDCNetworkList); i++ {
		if PDCNetworkList[i].Name == networkname {
			whattoreturn = PDCNetworkList[i]
		}

	}
	return (whattoreturn)
}

type NetworkSt struct {
	Name string
	VLANID string
	Gateway string
	Netmask string
	Bitmask string
	Location string
}

var PDCNetworkList []NetworkSt



