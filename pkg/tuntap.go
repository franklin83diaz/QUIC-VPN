package pkg

import (
	"fmt"
	"log"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// Create a new TUN device
func CreateVNet() {

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
	addr, err := netlink.ParseAddr("192.168.45.1/24")
	if err != nil {
		log.Fatal(err)
	}

	if err := netlink.AddrAdd(tun, addr); err != nil {
		log.Fatal(err)
	}

	packet := make([]byte, 9000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
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
