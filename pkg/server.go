package pkg

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"QUIC-VPN/utils"

	"github.com/quic-go/quic-go"
)

var (
	mutex sync.Mutex
)

func Server(tunFile *os.File) {

	// server address
	addr := "0.0.0.0:4242"

	// Create a QUIC listener
	listener, err := quic.ListenAddr(addr, utils.GenerateTLSConfigServer(), utils.GenerateQUICConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	defer tunFile.Close()
	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		go func() {

			stream, err := conn.AcceptStream(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			defer stream.Close()

			fmt.Println("New connection!")

			go func() {

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
			}()

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
		}()

	}

}
