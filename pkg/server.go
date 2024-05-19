package pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"QUIC-VPN/utils"

	"github.com/quic-go/quic-go"
)

func Server(tunFile *os.File) {

	// server address
	addr := "0.0.0.0:4242"
	ch := make(chan bool)

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

			go redirectTunToQuic(tunFile, stream)
			go redirectQuicToTun(stream, tunFile)

			<-ch
		}()

	}

}
