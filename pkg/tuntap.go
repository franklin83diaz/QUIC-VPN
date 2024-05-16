package pkg

import (
	"fmt"
	"log"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// Create a new TUN device
func CreateTun(ipServer string) {

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	tun := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: ifce.Name(),
			MTU:  1500,
		}}

	// Set the IP address and netmask for the TUN device
	if err := netlink.LinkSetUp(tun); err != nil {
		log.Fatal(err)
	}

	// Set the IP address and netmask for the TUN device
	SetIp(ipServer, ifce.Name())

	packet := make([]byte, 9000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
		log.Printf("%s", string(packet[:n]))

	}

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

// Delete the TUN device
func DeleteVNet() error {
	link, err := netlink.LinkByName("quic-tun0")
	if err != nil {
		fmt.Printf("Error to delete interface: %v\n", err)
		return err
	}

	// Delete the TUN device
	if err := netlink.LinkDel(link); err != nil {
		fmt.Printf("Error to delete interface: %v\n", err)
		return err
	}
	return nil
}
