package main

import (
	"fmt"
	"log"
	"net"
	"os"
)


func main(){

	hostname, err := os.Hostname()
	if err != nil{
		log.Fatalln("Failed to lookup hostname")
	}


	ip, err := net.LookupIP(hostname)
	if err != nil{
		log.Fatalln("Failed to look up ip")
	}

	fmt.Println(hostname)
	fmt.Println(ip)
	
	// loop for ever
	for{

	}
}
