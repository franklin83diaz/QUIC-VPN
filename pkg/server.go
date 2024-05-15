package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"QUIC-VPN/utils"

	"github.com/quic-go/quic-go"
)

func Server() {
	// server address
	addr := "127.0.0.1:4242"

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

				fmt.Printf("\033[32mServer: Received '%s'\033[0m\n", data)
				time.Sleep(1 * time.Second)
			}
		}()
	}

}
