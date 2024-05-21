package pkg

import (
	"QUIC-VPN/utils"
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func redirectTunToQuic(tunFile *os.File, stream []quic.Stream) {
	dataIn := make([]byte, 65536)
	for i := 0; i < 1; i++ {
		fmt.Println("stream ToQUIC: ", i)
		go func() {
			for {
				// Read from the TUN interface
				n, err := tunFile.Read(dataIn)
				if err != nil {
					log.Fatal(err)
				}

				// Send the data to the QUIC stream
				_, err = stream[i].Write(dataIn[:n])
				if err != nil {
					log.Fatal(err)
				}

			}
		}()
	}
}

func redirectQuicToTun(stream []quic.Stream, tunFile *os.File) {

	for i := 0; i < 1; i++ {

		fmt.Println("stream ToTun: ", i)

		dataOut := make([]byte, 65536)
		reader := bufio.NewReaderSize(stream[i], 10065536)
		totalLength := uint16(0)

		go func() {
			for {

				b, _ := reader.Peek(7)

				if string(b) == "initial" {
					reader.Discard(7)
					continue
				}

				//check if the packet is an IPv4 packet
				if !utils.IsIPv4(b) {
					totalLength6 := utils.GetTotalLengthIPv6(b) + uint16(40)
					reader.Discard(int(totalLength6))
					continue
				}

				readerLength := reader.Buffered()

				//TODO: Change to int
				totalLength = utils.GetTotalLength(b)

				if readerLength < int(totalLength) {
					//color red
					fmt.Println("\033[31m")
					fmt.Println("Buffered: ", readerLength)
					fmt.Println("Total: ", totalLength)
					fmt.Println("--------------------")
					//color reset
					fmt.Print("\033[0m")
					continue
				}

				bufferedLen := reader.Buffered()

				// // color red
				// fmt.Println("\033[31m")
				// fmt.Println("Buffered: ", bufferedLen)
				// fmt.Println("Total: ", totalLength)
				// fmt.Println("--------------------")
				// // color reset
				// fmt.Print("\033[0m")

				// Check if the buffer is less than the total length
				if uint16(bufferedLen) < totalLength {
					lengthToRead := uint16(bufferedLen)
					//bufferedLen = reader.Buffered()
					tt := 0
					//color green
					fmt.Println("\033[32m")
					fmt.Println("Buffered: ", bufferedLen)
					fmt.Println("\033[0m")

				read:
					// Read data from the QUIC stream
					//c: Change this for solve the problem when  missing data is too small and more than two packets
					readInt, err := reader.Read(dataOut[tt:(lengthToRead + uint16(tt))])
					if err != nil {
						log.Fatal(err)
					}

					// // color green
					// fmt.Println("\033[32m")
					// fmt.Println("lengthToRead: ", lengthToRead)
					// fmt.Println("\033[0m")

					// data missing
					lengthToRead = totalLength - uint16(lengthToRead)
					tt += readInt

					// fmt.Println("\033[32m")
					// fmt.Println("tt: ", tt, " < totalLength: ", totalLength)
					// fmt.Println("\033[0m")

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
		}()
	}
}
