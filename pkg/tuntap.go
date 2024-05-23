package pkg

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// Create a new TUN device
func CreateTun(cidr string) *os.File {

	fmt.Println("TUN interface created successfully")

	tun := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: "tun0",
			MTU:  1500,
		},
		Mode: netlink.TUNTAP_MODE_TUN,
	}

	// Add the interface to the system
	if err := netlink.LinkAdd(tun); err != nil {
		log.Fatalf("Failed to add the interface: %v", err)
	}

	//set txqlen
	if err := netlink.LinkSetTxQLen(tun, 1000); err != nil {
		log.Fatalf("Failed to set txqlen: %v", err)
	}

	// Disable allmulticast
	if err := netlink.LinkSetAllmulticastOff(tun); err != nil {
		log.Fatalf("Failed to set allmulticast off: %v", err)
	}

	// Bring the interface up
	if err := netlink.LinkSetUp(tun); err != nil {
		log.Fatalf("Failed to bring up the interface: %v", err)
	}

	// Set the IP address and netmask for the TUN device
	SetIp(cidr, "tun0")

	// open TUN
	tunFile, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		log.Panic("error open TUN: ", err)

	}

	// Configure TUN
	var ifr [unix.IFNAMSIZ + 64]byte
	copy(ifr[:unix.IFNAMSIZ], []byte("tun0\x00"))
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TUN | unix.IFF_NO_PI

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, tunFile.Fd(), uintptr(unix.TUNSETIFF), uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		log.Panic("rror configure TUN: ", errno)
	}

	return tunFile

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
