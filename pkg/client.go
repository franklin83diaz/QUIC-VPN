package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

// Client
func Client(ip string, port string, tunFile *os.File) {
	//cliente QUIC
	addr := ip + ":" + port

	// Create a QUIC session
	con, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), utils.GenerateQUICConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC stream
	stream, err := con.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer tunFile.Close()

	go func() {

		dataIn := make([]byte, 65536)

		for {
			// Read from the TUN interface
			n, err := tunFile.Read(dataIn)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				continue
			}

			err = utils.ValidateIPPacket(dataIn[:n])
			if err != nil {
				fmt.Println(err)
				continue
			}

			//Send the data to the QUIC stream
			_, err = stream.Write(dataIn[:n])
			if err != nil {
				log.Fatal(err)
			}

		}
	}()

	dataOut := make([]byte, 65536)

	for {

		// Read from the QUIC stream
		n, err := stream.Read(dataOut)
		if err != nil {
			log.Fatal(err)
		}

		// Write to the TUN interface
		_, err = tunFile.Write(dataOut[:n])
		if err != nil {
			fmt.Println(n)
			fmt.Println("-")
			//log.Fatal(err)
		}

	}

}
