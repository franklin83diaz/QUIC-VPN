package main

import (
	"QUIC-VPN/pkg"
	"time"
)

func main() {

	pkg.CreateVNet()
	time.Sleep(5 * time.Second)
	//pkg.DeleteVNet()

	go pkg.Server()
	time.Sleep(1 * time.Second)
	pkg.Client()

	time.Sleep(10 * time.Second)

}
