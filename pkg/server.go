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

				dataIn := make([]byte, 1500)

				for {
					// Read from the TUN interface
					n, err := tunFile.Read(dataIn)
					if err != nil {
						log.Fatal(err)
					}

					// Send the data to the QUIC stream
					_, err = stream.Write(dataIn[:n])
					if err != nil {
						log.Fatal(err)
					}

				}
			}()

			dataOut := make([]byte, 4096)
			reader := bufio.NewReader(stream)
			r := 0

			for {

				b, _ := reader.Peek(5)
				bufferedLen := reader.Buffered()
				totalLength := utils.GetToalLength(b)
				fmt.Println("Buffered Length: ", bufferedLen)
				fmt.Println("Total Length: ", totalLength)

				if totalLength == 0 {
					reader.Read(dataOut[:bufferedLen])
					bb, _ := reader.Peek(48)
					fmt.Println("Peek: ", bb)
					fmt.Println("Peek: ", string(bb))
					continue
				}

				if uint16(bufferedLen) < totalLength {
					if r < 5 {
						r++
						continue
					}
					r = 0

				}

				// Read from the QUIC stream
				n, err := reader.Read(dataOut[:totalLength])

				if err != nil {
					log.Fatal(err)
				}

				// Validar el paquete IP
				err = utils.ValidateIPPacket(dataOut[:n])
				if err != nil {
					fmt.Println(err)

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
