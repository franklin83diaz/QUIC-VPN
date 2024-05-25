package pkg

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TODO: ADD ICMP
// UT is UDP, TCP and IP together
// The head length is fixe to 11 Bytes

type UTProtocol struct {
	Length uint16 //2 bytes
	// 1 bytes 1 UDP, 21 TCP CONNECT, 22 TCP CONNECTED, 23 TCP DISCONNECT
	// The three-way handshake is done with the CONNECT in separate in the server and client
	// the four-way handshake is done with the DISCONNECT in separate in the server and client
	Protocol uint8
	IP       net.IP         // 4 bytes
	DsPort   layers.TCPPort //2 bytes
	SrcPort  layers.TCPPort //2 bytes
}

func UtToTcpIp(srcIP net.IP, destIP net.IP, srcPort layers.TCPPort, dstPort layers.TCPPort, sequenceNumber uint32, acknowledgmentNumber uint32, window uint16, syn bool, ack bool, fin bool) []byte {

	// Create IP layers
	ipLayer := &layers.IPv4{
		SrcIP:    srcIP,                // IP  source net.IP{192, 168, 1, 2}
		DstIP:    destIP,               // IP destination net.IP{192, 168, 1, 1}
		Protocol: layers.IPProtocolTCP, // Set Protocol TCP
	}

	// Create TCP layers
	// https://en.wikipedia.org/wiki/Transmission_Control_Protocol
	tcpLayer := &layers.TCP{
		SrcPort: srcPort,              // source port layers.TCPPort(54321)
		DstPort: dstPort,              // destination port layers.TCPPort(54321)
		SYN:     syn,                  // Flag SYN  start connection
		ACK:     ack,                  // Flag ACK  Acknowledgment of receipt
		FIN:     fin,                  // Flag FIN  end connection
		Ack:     acknowledgmentNumber, // other part connection sequence number + 1
		Seq:     sequenceNumber,       // initial sequence number
		Window:  window,               // size windows
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
