package pkg

import (
	"QUIC-VPN/utils"
	"encoding/binary"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func redirectTunToQuic(tunFile *os.File, stream quic.Stream) {

	dataIn := make([]byte, 65002)

	for {
		// Read from the TUN interface
		n, err := tunFile.Read(dataIn[2:])
		if err != nil {
			log.Fatal(err)
		}

		if !utils.IsIPv4(dataIn[2:]) {
			continue
		}

		binary.BigEndian.PutUint16(dataIn[0:2], uint16(n))

		dataOut := make([]byte, n+2)
		copy(dataOut, dataIn[:n+2])

		go func(d []byte) {
			// Send the data to the QUIC stream
			_, err = stream.Write(d)
			if err != nil {
				log.Fatal(err)
			}
		}(dataOut)

	}
}

func redirectQuicToTun(stream quic.Stream, tunFile *os.File) {
	data := make([]byte, 65000)

	for {
		lenChunk := make([]byte, 2)
		stream.Read(lenChunk)

		lenChunkInt := int(binary.BigEndian.Uint16(lenChunk))

		i, _ := stream.Read(data[:lenChunkInt])

		for i < lenChunkInt {
			ii, _ := stream.Read(data[i:lenChunkInt])
			i += ii
		}

		dataOut := make([]byte, lenChunkInt)
		copy(dataOut, data[:lenChunkInt])

		//	go func(d []byte) {

		// Write to the TUN interface
		_, err := tunFile.Write(dataOut)
		if err != nil {
			log.Fatal(err)
		}
		//	}(dataOut)
		//	time.Sleep(50 * time.Microsecond)

	}

}
