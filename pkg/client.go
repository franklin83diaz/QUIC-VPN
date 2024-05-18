package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"log"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/songgao/water"
)

// Client
func Client(ip string, port string, ifce *water.Interface) {
	//cliente QUIC
	addr := ip + ":" + port

	// Configuración QUIC
	quicConfig := &quic.Config{
		MaxIdleTimeout:       60 * time.Second,
		HandshakeIdleTimeout: 60 * time.Second,
	}

	// Create a QUIC session
	con, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), quicConfig)
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

		go func() {

			for n, _ := stream.Read(dataIn); n > 0; n, _ = stream.Read(dataIn) {
				_, err := ifce.Write(dataIn[:n])
				if err != nil {
					log.Println(err)
				}

			}
		}()
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
