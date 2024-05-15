package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"fmt"
	"log"
	"time"

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

	con2, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC stream
	stream, err := con.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	stream2, err := con2.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stream.Write([]byte("Hello, World!"))
	time.Sleep(1 * time.Second)
	stream2.Write([]byte("Hello, World! from 2nd stream"))
	time.Sleep(1 * time.Second)
	stream.Write([]byte("2nd message"))
	time.Sleep(1 * time.Second)
	stream.Write([]byte("3rd message"))
	stream.Close()
	time.Sleep(1 * time.Second)

	fmt.Println("\033[31mConnection closed\033[0m")

}
