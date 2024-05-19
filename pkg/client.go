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
	con, err := quic.DialAddr(context.Background(), addr, utils.GenerateTLSConfigClient(), utils.GenerateQUICConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Create a QUIC stream
	stream, err := con.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer tunFile.Close()
	defer stream.Close()

	go redirectTunToQuic(tunFile, stream)
	go redirectQuicToTun(stream, tunFile)
	ch <- true

}
