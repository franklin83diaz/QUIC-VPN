package pkg

import (
	"QUIC-VPN/utils"
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func redirectTunToQuic(tunFile *os.File, stream quic.Stream) {

	dataIn := make([]byte, 65536)

	for {
		// Read from the TUN interface
		n, err := tunFile.Read(dataIn)
		if err != nil {
			log.Fatal(err)
		}

		// check if the packet is valid
		err = utils.ValidateIPPacket(dataIn[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Send the data to the QUIC stream
		_, err = stream.Write(dataIn[:n])
		if err != nil {
			log.Fatal(err)
		}

	}
}

func redirectQuicToTun(stream quic.Stream, tunFile *os.File) {
	dataOut := make([]byte, 65536)
	reader := bufio.NewReaderSize(stream, 65536)

	for {

		b, _ := reader.Peek(5)
		bufferedLen := reader.Buffered()
		totalLength := utils.GetTotalLength(b)

		// Check if the buffer is less than the total length
		if bufferedLen < totalLength {
			lengthToRead := bufferedLen
			tt := 0

		read:
			// Read data from the QUIC stream
			//c: Change this for solve the problem when  missing data is too small and more than two packets
			readInt, err := reader.Read(dataOut[tt:(lengthToRead + tt)])
			if err != nil {
				log.Fatal(err)
			}

			// data missing
			lengthToRead = totalLength - lengthToRead
			tt += readInt
			if tt < int(totalLength) {
				goto read
			}

		} else {
			// Read from the QUIC stream
			_, err := reader.Read(dataOut[:totalLength])
			if err != nil {
				log.Fatal(err)
			}
		}

		err := utils.ValidateIPPacket(dataOut[:totalLength])
		if err != nil {
			reader.Discard(totalLength)
			continue
		}

		copyDataOut := make([]byte, totalLength)
		copy(copyDataOut, dataOut[:totalLength])

		go func(data []byte) {

			// Write to the TUN interface
			_, err = tunFile.Write(copyDataOut)
			if err != nil {
				log.Println(err)
			}
		}(copyDataOut)

	}
}
