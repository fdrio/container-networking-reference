package main

import (
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/vishvananda/netlink"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP         string `json:"ip"`
}

func cmdAdd(args *skel.CmdArgs) error {
	sb := SimpleBridge{}
	if err := json.Unmarshal(args.StdinData, &sb); err != nil{
		return err
	}
	fmt.Println(sb)


	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: sb.BridgeName,
			MTU: 1500, // (not including a 4 byte header)

			// Let kernel use default txqueuelen; leaving it unset
            // means 0, and a zero-length TX queue messes up FIFO
            // traffic shapers which use TX queue length as the
            // default packet limit
			TxQLen: -1, // Le the kernel decide by itself. It knows best.
		},
	}

	err := netlink.LinkAdd(br)

	if err!=nil &&  err != syscall.EEXIST{
		return err
	}
	
	// Create netlink ipAddr from input IP
	/*ipAddr, err := netlink.ParseAddr(sb.IP)
	if err != nil{
		log.Println("Error parsing ip")
		log.Println(err)
	}

	// Add the ip to the device (Same as ip addr add $BRIDGE_IP) default subnet is /24
	netlink.AddrAdd(br, ipAddr)
	*/

	// same as ip link set $BRIDGE_NAME up
	if err := netlink.LinkSetUp(br); err != nil{
		return err
	}
	return nil
}

func cmdCheck(args *skel.CmdArgs) error {
	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "Hello World CNI")
}
