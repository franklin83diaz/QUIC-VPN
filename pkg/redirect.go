package pkg

import (
	"QUIC-VPN/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

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
	totalLength := uint16(0)

	for {

		b, _ := reader.Peek(5)
		bufferedLen := reader.Buffered()
		totalLength = utils.GetTotalLength(b)

		// color red
		fmt.Println("\033[31m")
		fmt.Println("Buffered: ", bufferedLen)
		fmt.Println("Total: ", totalLength)
		fmt.Println("--------------------")
		// color reset
		fmt.Print("\033[0m")

		if uint16(bufferedLen) < totalLength {
			time.Sleep(500 * time.Millisecond)
			bufferedLen = reader.Buffered()
			fmt.Println("Buffered: ", bufferedLen)
			//TODO: Fix this
			log.Println("The problem is that the buffer does not enter data again until it is read.")
			continue
		}

		// Read from the QUIC stream
		n, err := reader.Read(dataOut[:totalLength])
		if err != nil {
			log.Fatal(err)
		}

		err = utils.ValidateIPPacket(dataOut[:n])
		if err != nil {
			fmt.Println(err)
			fmt.Println("packet length: ", totalLength)
			fmt.Println("data length: ", n)
			continue
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
