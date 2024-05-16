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

			data := make([]byte, 1500)
			for n, _ := stream.Read(data); n > 0; n, _ = stream.Read(data) {
				ifce.Write(data[:n])
			}

			packet := make([]byte, 1500)
			for {
				n, err := ifce.Read(packet)
				if err != nil {
					log.Fatal(err)
				}
				_, err = stream.Write(packet[:n])
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

	}

}
