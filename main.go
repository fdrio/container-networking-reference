/*package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

type NetInfo struct {
	BridgeName string `json:"bridgeName"`
	BridgeIP   string `json:"bridgeIP"`
	VethConIP  string `json:"vethConIP"`
}

func init() {
	runtime.LockOSThread()
}

func cmdAdd(args *skel.CmdArgs) error {
	netInfo := NetInfo{}
	if err := json.Unmarshal(args.StdinData, &netInfo); err != nil {
		return err
	}
	fmt.Println(netInfo)
/////////////////////////////////// Add Bridge /////////////////////////////////// 
	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: netInfo.BridgeName,
			MTU:  1500, // (not including a 4 byte header)

			// Let kernel use default txqueuelen; leaving it unset
			// means 0, and a zero-length TX queue messes up FIFO
			// traffic shapers which use TX queue length as the
			// default packet limit
			TxQLen: -1, // Le the kernel decide by itself. It knows best.
		},
	}

	err := netlink.LinkAdd(br)
	if err != nil && err != syscall.EEXIST {
		log.Println("Error adding new bridge")
		//Check error and whether the link exist already
	}
////////////////////////////////////////////////////////////////////// 
	//Get link by name

	l, err := netlink.LinkByName(netInfo.BridgeName)

	if err != nil {
		log.Println("Error finding device by name")
		return err
	}

/////////////////////////////////// Link to Bridge  & add ip ///////////////////////////////////  

// Make sure the link is of type bridge

	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		log.Println("Link name already exists and is of another type")
		return err
	}

	
///////////////////////////////////////////////////////////////////////////////////////////////  


//Get the network namespace path from args
	networkNamespace, err := ns.GetNS(args.Netns)

	if err != nil {
		log.Println("Error getting namespace from argument path")
		return err
	}

	hostIface := &current.Interface{}

	var handler = func(hostns ns.NetNS) error {

			//Create veth pair
			hostVeth, containerVeth, err := ip.SetupVeth(args.IfName, 1500, hostns)
			if err != nil {
				log.Println("Error creating veth pair")
				return err
			}

			hostIface.Name = hostVeth.Name
			// Parse CIDR and assign ip
			conVethLink, err := netlink.LinkByName(containerVeth.Name)
			if err != nil {
				log.Println("Error finding container veth by name")
				return err
			}

			addr, err := netlink.ParseAddr(netInfo.VethConIP)
			if err != nil {
				log.Println("Error parsing container veth name ip address")
				return err
			}
	
			
			// add ipaddr to container veth end of the link
			if err = netlink.AddrAdd(conVethLink, addr); err != nil {
				log.Println("Error adding address to link")
				return err
			}

			//Bring containerVeth up
			if err = netlink.LinkSetUp(conVethLink); err != nil {
				log.Println("Error bringin container veth up")
			}
			// Bring Network Namespace loopback up
			lo, err := netlink.LinkByName("lo")
			if err != nil{
				log.Println("Error finding loopback in netns")
				return err
			}
			if err = netlink.LinkSetUp(lo); err != nil{
				log.Println("Error bringing up loopback interface")
				return err
			}
		
			// Get bridge ipnet 
			brnet, err := netlink.ParseIPNet(netInfo.BridgeIP)
			if err != nil{
				log.Println("Error parsing bridge ip net")
				return err
			}
			// Add route from inside network namespace
			route := netlink.Route{
				LinkIndex: conVethLink.Attrs().Index, 
				Src: addr.IP,
				Dst: brnet,	
			}

			if err := netlink.RouteAdd(&route); err != nil{
				log.Println("Error adding route")
				return err
			}
			
		return nil
	}

	if err := networkNamespace.Do(handler); err != nil {
		log.Println("Error applying function in network namespace")
		return err
	}

	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		log.Println("Error finding host veth by name")
		return err
	}

	//bring the veth end of the pair that is attached to the bridge up
	if err := netlink.LinkSetUp(hostVeth); err != nil{
		log.Println("Error bringing host veth up")
	}
	// Setting bridge as master 

	if err := netlink.LinkSetMaster(hostVeth, newBr); err != nil {
		log.Println("Error setting up host veth to master bridge")
		return err
	}

	// same as ip link set $BRIDGE_NAME up
	if err := netlink.LinkSetUp(newBr); err != nil {
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
*/
