package pkg

import (
	"log"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// Create a new TUN device
func CreateTun(ipServer string) *water.Interface {

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	tun := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifce.Name(),
			MTU:  1500,
		}}

	// Set the IP address and netmask for the TUN device
	if err := netlink.LinkSetUp(tun); err != nil {
		log.Fatalln(err)
	}

	// Set the IP address and netmask for the TUN device
	SetIp(ipServer, ifce.Name())

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
