package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ifaces, err := net.Interfaces()

	if err != nil {
		log.Fatalln("Error getting interfaces")

	}

	ifaceMap := make(map[string]string)
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()

		if err != nil {
			log.Fatalln(err)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.To4() != nil {
				ifaceMap[iface.Name]=ip.String()
				fmt.Printf("Name: %s \n Addr: %s \n", iface.Name, ip.String())
			}

		}
	}


	fmt.Println(ifaceMap)

	for {
	}
}
