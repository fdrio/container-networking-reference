package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP         string `json:"ip"`
}

func init(){
	runtime.LockOSThread()
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
	
	// same as ip link set $BRIDGE_NAME up
	if err := netlink.LinkSetUp(br); err != nil{
		log.Println("Error bringing up interface")
		return err
	}

	l, err := netlink.LinkByName(sb.BridgeName)
	if err != nil{
		log.Println("Error finding link by name")
		return err
	}
	// Make sure the link is of type bridge (netlink.Bridge)
	newBr, ok := l.(*netlink.Bridge)

	if !ok{
		return fmt.Errorf("%q already exists but is not a bridge", sb.BridgeName)
	}

	netns, err := ns.GetNS(args.Netns)
	
	if err != nil{
		return err
	}

	hostIface := &current.Interface{}

	// this handler func creates a vethpair and we get the name of the host side veth
	var handler = func(hostNS ns.NetNS)error{
		//hostVeth, containerVeth, err := ip.SetupVeth(args.IfName, 1500, hostNS)
		// Creates a veth pair, sets mtu and moves one side of the veth pair to hostNS
		hostVeth, _, err := ip.SetupVeth(args.IfName, 1500, hostNS)
		if err != nil{
			return err
		}
	
		//Get the name of the host side of veth pair
		hostIface.Name = hostVeth.Name

		ipv4Addr, ipv4Net, err := net.ParseCIDR(sb.IP)
		if err != nil{
			log.Println("Error parsing ip address")
		}
		ipv4Net.IP = ipv4Addr
		addr:=&netlink.Addr{IPNet: ipv4Net, Label:""}
		// assign the address to the bridge
		if err := netlink.AddrAdd(newBr,addr); err != nil{
			return err
		}
		return nil
	}


	if err := netns.Do(handler); err != nil{
		return err
	}



	hostVeth, err := netlink.LinkByName(hostIface.Name)

	if err != nil{
		return err
	}

	if err := netlink.LinkSetMaster(hostVeth, newBr); err != nil{
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
