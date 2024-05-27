package main

import (
	"QUIC-VPN/pkg"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	f, _ := os.Create("profile-client.pprof")
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ip := os.Args[1]

	// Create netowrk interface
	osFile := pkg.CreateTun("192.168.45.254/24")

	pkg.Client(ip, "4242", osFile)

	time.Sleep(10 * time.Second)

}
