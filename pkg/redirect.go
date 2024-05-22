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

		totalLength := utils.GetTotalLength(dataIn[:n])
		datalen := len(dataIn[:n])

		// check if the packet is valid
		err = utils.ValidateIPPacket(dataIn[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("--------read-tun------------")
		fmt.Println("Total Length: ", totalLength)
		fmt.Println("Data Length: ", datalen)
		fmt.Println("--------------------")

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

		// color red
		fmt.Println("\033[31m")
		fmt.Println("Buffered: ", bufferedLen)
		fmt.Println("Total: ", totalLength)
		fmt.Println("--------------------")
		// color reset
		fmt.Print("\033[0m")

		// Check if the buffer is less than the total length
		if bufferedLen < totalLength {
			lengthToRead := bufferedLen
			bufferedLen = reader.Buffered()
			tt := 0
			// color green
			fmt.Println("\033[32m")
			fmt.Println("Buffered: ", bufferedLen)
			fmt.Println("\033[0m")

		read:
			// Read data from the QUIC stream
			//c: Change this for solve the problem when  missing data is too small and more than two packets
			readInt, err := reader.Read(dataOut[tt:(lengthToRead + tt)])
			if err != nil {
				log.Fatal(err)
			}

			// color green
			fmt.Println("\033[32m")
			fmt.Println("lengthToRead: ", lengthToRead)
			fmt.Println("\033[0m")

			// data missing
			lengthToRead = totalLength - lengthToRead
			tt += readInt
			fmt.Println("\033[32m")
			fmt.Println("tt: ", tt, " < totalLength: ", totalLength)
			fmt.Println("\033[0m")
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
			//color yellow
			fmt.Println("\033[33m")
			fmt.Println(err)
			fmt.Println("packet length: ", totalLength)
			fmt.Println("data length: ", totalLength)
			fmt.Println("\033[0m")
			continue
		}

		// Write to the TUN interface
		_, err = tunFile.Write(dataOut[:totalLength])
		if err != nil {
			fmt.Println(totalLength)
			fmt.Println("-")
			fmt.Println(err)

		}

	}
}
