package pkg

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// UT is UDP, TCP and IP together
// The head length is fixe to 11 Bytes
type UTProtocol struct {
	Length   uint16         //2 bytes
	Protocol uint8          // 1 bytes 1 UDP, 21 SYN, 22 SYN-ACK, 23 ACK, 24 FIN
	IP       net.IP         // 4 bytes
	DsPort   layers.TCPPort //2 bytes
	SrcPort  layers.TCPPort //2 bytes
}

func UtToTcpIp(ut UTProtocol, destIP net.IP, sequenceNumber uint32, acknowledgmentNumber uint32, window uint16) []byte {

	// Create IP layers
	ipLayer := &layers.IPv4{
		SrcIP:    ut.IP,                // IP  source net.IP{192, 168, 1, 2}
		DstIP:    destIP,               // IP destination net.IP{192, 168, 1, 1}
		Protocol: layers.IPProtocolTCP, // Set Protocol TCP
	}

	// Create TCP layers
	// https://en.wikipedia.org/wiki/Transmission_Control_Protocol
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(80),    // source port
		DstPort: layers.TCPPort(54321), // destination port
		SYN:     true,                  // Flag SYN  start connection
		ACK:     true,                  // Flag ACK  Acknowledgment of receipt
		Ack:     acknowledgmentNumber,  // other part connection sequence number + 1
		Seq:     sequenceNumber,        // initial sequence number
		Window:  window,                // size windows
	}

	// Is Important call SetNetworkLayerForChecksum to let gopacket know
	// how to calculate the TCP checksum over the IPv4 packet
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	// Serialize the layers to build the packet
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	// Serialize the packet
	err := gopacket.SerializeLayers(buffer, options, ipLayer, tcpLayer)
	if err != nil {
		log.Fatal(err)
	}
	headTcpIp := buffer.Bytes()

	return headTcpIp
}
