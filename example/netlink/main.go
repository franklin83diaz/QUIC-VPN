package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func main() {
	// Crear una nueva interfaz de tipo TUN
	link := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: "tun0",
		},
		Mode: netlink.TUNTAP_MODE_TUN,
	}

	// Agregar la interfaz al sistema
	if err := netlink.LinkAdd(link); err != nil {
		fmt.Printf("Error al agregar la interfaz: %v\n", err)
		return
	}

	// Traer la interfaz arriba (up)
	if err := netlink.LinkSetUp(link); err != nil {
		fmt.Printf("Error al levantar la interfaz: %v\n", err)
		return
	}

	fmt.Println("Interfaz TUN tun0 creada y levantada con éxito")

	// Abrir el dispositivo TUN
	tunFile, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("Error al abrir el dispositivo TUN: %v\n", err)
		return
	}
	defer tunFile.Close()

	// Configurar el dispositivo TUN
	var ifr [unix.IFNAMSIZ + 64]byte
	copy(ifr[:unix.IFNAMSIZ], []byte("tun0\x00"))
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TUN | unix.IFF_NO_PI

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, tunFile.Fd(), uintptr(unix.TUNSETIFF), uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		fmt.Printf("Error al configurar el dispositivo TUN: %v\n", errno)
		return
	}

	// Escribir datos en la interfaz TUN
	data := []byte("Hello, TUN interface!")
	if _, err := tunFile.Write(data); err != nil {
		fmt.Printf("Error al escribir en la interfaz TUN: %v\n", err)
		return
	}

	fmt.Println("Datos escritos en la interfaz TUN con éxito")

	// Leer datos desde la interfaz TUN
	buf := make([]byte, 1500) // Tamaño típico de MTU
	n, err := tunFile.Read(buf)
	if err != nil {
		fmt.Printf("Error al leer desde la interfaz TUN: %v\n", err)
		return
	}

	fmt.Printf("Datos leídos desde la interfaz TUN: %s\n", string(buf[:n]))
}
