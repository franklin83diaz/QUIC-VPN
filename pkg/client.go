package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

// Client
func Client(ip string, port string, tunFile *os.File) {
	//cliente QUIC
	addr := ip + ":" + port
	ch := make(chan bool)

	// Create a QUIC session
	conn, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), utils.GenerateQUICConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC stream

	if err != nil {
		log.Fatal(err)
	}

	defer tunFile.Close()

	go redirectTunToQuic(tunFile, conn)
	go redirectQuicToTun(conn, tunFile)
	ch <- true

}
