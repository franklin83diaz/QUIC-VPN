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

	streams := []quic.Stream{}

	for i := 0; i < 1; i++ {
		// Create a QUIC stream
		stream, err := con.OpenStream()
		if err != nil {
			log.Fatal(err)
		}
		// Send initial packet
		stream.Write([]byte("initial"))
		streams = append(streams, stream)
	}

	defer tunFile.Close()
	defer func() {
		for _, stream := range streams {
			stream.Close()
		}
	}()

	go redirectTunToQuic(tunFile, streams)
	go redirectQuicToTun(streams, tunFile)
	ch <- true

}
