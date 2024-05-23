package pkg

import (
	"QUIC-VPN/utils"
	"encoding/binary"
	"fmt"
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

		fmt.Println("\033[32m", "n: ", n, "\033[0m")

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
	dataOut := make([]byte, 1500)
	preStreamLen := 0

	for {
		temp := make([]byte, 1500)
		packetVpnLen := make([]byte, 2)
		_, err := stream.Read(packetVpnLen)
		if err != nil {
			log.Fatal(err)
		}
		packetVpnLenInt := int(binary.BigEndian.Uint16(packetVpnLen))
		fmt.Println("packetVpnLenInt: ", packetVpnLenInt)

		streamLen, _ := stream.Read(temp)

		if streamLen < packetVpnLenInt {
			copy(dataOut, temp[preStreamLen:streamLen+preStreamLen])
			preStreamLen = streamLen
			continue
		} else {
			copy(dataOut, temp[preStreamLen:streamLen+preStreamLen])
		}
		preStreamLen = 0

		err = utils.ValidateIPPacket(dataOut[:packetVpnLenInt])
		if err != nil {
			continue
		}

		copyDataOut := make([]byte, packetVpnLenInt)
		copy(copyDataOut, dataOut[:packetVpnLenInt])

		go func(data []byte) {

			// Write to the TUN interface
			_, err = tunFile.Write(copyDataOut)
			if err != nil {
				log.Println(err)
			}
		}(copyDataOut)

	}
}
