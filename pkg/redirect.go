package pkg

import (
	"QUIC-VPN/utils"
	"encoding/binary"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func redirectTunToQuic(tunFile *os.File, stream quic.Stream) {

	dataIn := make([]byte, 1502)

	for {
		// Read from the TUN interface
		n, err := tunFile.Read(dataIn[2:])
		if err != nil {
			log.Fatal(err)
		}

		if !utils.IsIPv4(dataIn[2:]) {
			continue
		}

		//fmt.Println("\033[32m", "n: ", n, "\033[0m")

		binary.BigEndian.PutUint16(dataIn[0:2], uint16(n))

		//TODO:
		// implement other bit for check is the number a diferent bwteen dataIn[0] and dataIn[1]
		// add consecutive for now where restaer and add cache for data

		// Send the data to the QUIC stream
		_, err = stream.Write(dataIn[:n+2])
		if err != nil {
			log.Fatal(err)
		}

	}
}

func redirectQuicToTun(stream quic.Stream, tunFile *os.File) {
	data := make([]byte, 1500)

	for {
		lenchunk := make([]byte, 2)
		stream.Read(lenchunk)

		lenchunkInt := int(binary.BigEndian.Uint16(lenchunk))

		i, _ := stream.Read(data[:lenchunkInt])

		for i < lenchunkInt {
			ii, _ := stream.Read(data[i:lenchunkInt])
			i += ii
		}

		dataOut := make([]byte, lenchunkInt)
		copy(dataOut, data[:lenchunkInt])

		//go func(d []byte) {

		// Write to the TUN interface
		_, err := tunFile.Write(dataOut)
		if err != nil {
			log.Fatal(err)
		}
		//}(dataOut)

	}
}
