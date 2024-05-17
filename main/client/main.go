package main

import (
	"QUIC-VPN/pkg"
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var ip string
	//ask for the server
	fmt.Print("Ip del servidor: ")
	fmt.Scanln(&ip)

	// Create netowrk interface
	ifce := pkg.CreateTun("192.168.45.254/24")

	pkg.Client(ip, "4242", ifce)

	time.Sleep(10 * time.Second)

}
