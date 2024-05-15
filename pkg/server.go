package pkg

import (
	"context"
	"fmt"
	"log"

	"QUIC-VPN/utils"

	"github.com/quic-go/quic-go"
)

func Server() {
	// server address
	addr := ":4242"

	// Create a QUIC listener
	listener, err := quic.ListenAddr(addr, utils.GenerateTLSConfigServer(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Accept QUIC connections
	for {
		con, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection accepted")
		// Handle the session
		go handleSession(con)
	}
}

func handleSession(con quic.Connection) {
	// Accept streams
	for {
		stream, err := con.AcceptStream(context.Background())
		if err != nil {
			log.Println(err)
			return
		}

		// Handle the stream
		go handleStream(stream)
	}
}

func handleStream(stream quic.Stream) {
	// Read from the stream
	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	// Print the message
	log.Printf("Received: %s", buf[:n])

	// Write to the stream
	_, err = stream.Write([]byte("Hello, World!"))
	if err != nil {
		log.Println(err)
		return
	}
}
