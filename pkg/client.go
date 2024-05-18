package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"log"

	"github.com/quic-go/quic-go"
	"github.com/songgao/water"
)

// Client
func Client(ip string, port string, ifce *water.Interface) {
	//cliente QUIC
	addr := ip + ":" + port

	// Create a QUIC session
	con, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC stream
	stream, err := con.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		dataIn := make([]byte, 1500)

		for {
			// Read from the TUN interface
			n, err := ifce.Read(dataIn)
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

	dataOut := make([]byte, 1500)

	for {
		// Read from the QUIC stream
		n, err := stream.Read(dataOut)
		if err != nil {
			log.Fatal(err)
		}

		// Write to the TUN interface
		_, err = ifce.Write(dataOut[:n])
		if err != nil {
			log.Fatal(err)
		}

	}

}
