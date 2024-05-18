package pkg

import (
	"fmt"
	"log"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// Create a new TUN device
func CreateTun(cidr string) *water.Interface {

	fmt.Println("TUN interface created successfully")

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			MultiQueue: true,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	tun := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifce.Name(),
		},
		Mode: netlink.TUNTAP_MODE_TUN,
	}

	// Change MTU
	// if err := netlink.LinkSetMTU(tun, 1350); err != nil {
	// 	log.Fatalf("Failed to set MTU: %v", err)
	// }

	// Qlen to 1000
	if err := netlink.LinkSetTxQLen(tun, 100); err != nil {
		log.Fatalf("Failed to set TxQLen: %v", err)
	}

	// Bring the interface up
	if err := netlink.LinkSetUp(tun); err != nil {
		log.Fatalf("Failed to bring up the interface: %v", err)
	}

	// Set the IP address and netmask for the TUN device
	SetIp(cidr, ifce.Name())

	return ifce

}

func SetIp(ip string, ifceName string) {
	// Set the IP address and netmask for the TUN device
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		log.Fatal(err)
	}

	link, err := netlink.LinkByName(ifceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		log.Fatal(err)
	}
}
