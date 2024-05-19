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

		ok := utils.IsIPv4(dataIn[:n])
		if !ok {
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
	r := 0
	max := 0

	for {
		fmt.Println("Max: ", max)

		b, _ := reader.Peek(5)
		bufferedLen := reader.Buffered()
		totalLength := utils.GetTotalLength(b)
		// fmt.Println("Buffered Length: ", bufferedLen)
		// fmt.Println("Total Length: ", totalLength)

		if totalLength > uint16(max) {
			max = int(totalLength)
		}

		if totalLength == 0 {
			reader.Read(dataOut[:bufferedLen])
			bb, _ := reader.Peek(48)
			fmt.Println("Peek: ", bb)
			fmt.Println("Peek: ", string(bb))
			continue
		}

		if uint16(bufferedLen) < totalLength {
			if r < 10 {
				r++
				continue
			}
			r = 0
			totalLength = uint16(bufferedLen)

		}

		// Read from the QUIC stream
		n, err := reader.Read(dataOut[:totalLength])

		if err != nil {
			log.Fatal(err)
		}

		// Validar el paquete IP
		err = utils.ValidateIPPacket(dataOut[:n])
		if err != nil {
			continue
			//fmt.Println(err)
		}

		// Write to the TUN interface
		_, err = tunFile.Write(dataOut[:n])
		if err != nil {

			fmt.Println(n)
			fmt.Println("-")
			fmt.Println(err)

		}

	}
}
