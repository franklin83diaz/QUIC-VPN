package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"fmt"
	"log"

	"github.com/quic-go/quic-go"
)

// Client
func Client() {
	//cliente QUIC
	addr := "127.0.0.1:4242"

	// Create a QUIC session
	con, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to", addr)

	// Create a QUIC stream
	stream, err := con.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stream.Write([]byte("Hello, World!"))
	stream.Close()

}
