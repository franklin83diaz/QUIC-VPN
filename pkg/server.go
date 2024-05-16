package pkg

import (
	"context"
	"log"

	"QUIC-VPN/utils"

	"github.com/quic-go/quic-go"
	"github.com/songgao/water"
)

func Server(ifce *water.Interface) {

	// server address
	addr := "0.0.0.0:4242"

	// Create a QUIC listener
	listener, err := quic.ListenAddr(addr, utils.GenerateTLSConfigServer(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
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

			// Read data from the stream and write it to the tun interface
			go func() {
				dataIn := make([]byte, 1500)
				for n, _ := stream.Read(dataIn); n > 0; n, _ = stream.Read(dataIn) {
					ifce.Write(dataIn[:n])
				}
			}()

			// Read data from the tun interface and write it to the stream
			dataOut := make([]byte, 1500)
			for {
				n, err := ifce.Read(dataOut)
				if err != nil {
					log.Fatal(err)
				}
				_, err = stream.Write(dataOut[:n])
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

	}

}
