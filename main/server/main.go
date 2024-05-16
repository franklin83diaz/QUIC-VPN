package main

import (
	"QUIC-VPN/internal"
	"QUIC-VPN/pkg"
	"QUIC-VPN/utils"
	"fmt"
	"strings"
)

func main() {
	// Get database connection
	db := pkg.Db

	// Find network
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

	ipServer = ipServer + "/" + strings.Split(network.CIDR, "/")[1]

	fmt.Println("Network CIDR: ", network.CIDR)

	// Create Virtual Network Interface
	ifce := pkg.CreateTun(ipServer)

	// Connection QUIC listener
	pkg.Server(ifce)

	// Create user admin for debug
	// GODEBUG := os.Getenv("GODEBUG")
	// if GODEBUG == "1" {

	// 	user := pkg.User{}
	// 	db.AutoMigrate(&user)

	// 	tx := db.Where("username = ?", "admin").First(&user)
	// 	if tx.Error != nil {
	// 		if tx.Error == gorm.ErrRecordNotFound {
	// 			//create user
	// 			user = pkg.User{
	// 				Name:         "admin",
	// 				Last:         "admin",
	// 				Username:     "admin",
	// 				HashPassword: utils.HashPassword("admin"),
	// 				Email:        "admin@noemail.com",
	// 				Active:       true,
	// 				Ip:           "192.168.145.254",
	// 			}

	// 			db.Create(&user)
	// 		} else {
	// 			fmt.Printf("Error : %v\n", tx.Error)
	// 		}
	// 	}
	// }

}
