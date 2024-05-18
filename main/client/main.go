package main

import (
	"QUIC-VPN/pkg"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ip := os.Args[1]

	// Create netowrk interface
	osFile := pkg.CreateTun("192.168.45.254/24")

	pkg.Client(ip, "4242", osFile)

	time.Sleep(10 * time.Second)

}
