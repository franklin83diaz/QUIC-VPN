package main

import (
	"QUIC-VPN/pkg"
	"time"
)

func main() {
	go pkg.Server()
	time.Sleep(1 * time.Second)
	pkg.Client()
}
