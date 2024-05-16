package main

import (
	"QUIC-VPN/internal"
	"QUIC-VPN/pkg"
	"QUIC-VPN/utils"
	"fmt"
	"strings"
)

func main() {

	// Find network
	db := pkg.Db
	network := pkg.Network{}

	db.AutoMigrate(&network)

	db.First(&network, 1)

	if network.CIDR == "" {
		network.CIDR = internal.DefaultCIDR
		network.ID = 1
		db.Create(&network)
	}

	_, ipServer, _, err := utils.AvailableIPs(network.CIDR)
	if err != nil {
		fmt.Println(err)
	}

	ipServer = ipServer + strings.Split(network.CIDR, "/")[1]

	fmt.Println("Network CIDR: ", network.CIDR)

	// Create Virtual Network Interface
	pkg.CreateTun(ipServer)

}
