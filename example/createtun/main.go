package main

import (
	"QUIC-VPN/pkg"
	"log"
	"time"
)

func main() {

	//log flags show file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pkg.CreateTun("192.168.10.1/24")
	time.Sleep(1000 * time.Second)

}
