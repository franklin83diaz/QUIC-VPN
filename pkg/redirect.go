package pkg

import (
	"QUIC-VPN/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func redirectTunToQuic(tunFile *os.File, conn quic.Connection) {
	dataIn := make([]byte, 6500)

	for {
		// Read from the TUN interface
		n, err := tunFile.Read(dataIn)
		if err != nil {
			log.Fatal(err)
		}

		// check if the packet is valid
		err = utils.ValidateIPPacket(dataIn[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}

		go func() {
			stream, _ := conn.OpenUniStreamSync(context.Background())
			stream.Write(dataIn[:n])
			stream.Close()
		}()

	}
}

func redirectQuicToTun(conn quic.Connection, tunFile *os.File) {

	for {

		stream, _ := conn.AcceptUniStream(context.Background())
		go func() {
			dataOut := make([]byte, 65000)
			i, _ := stream.Read(dataOut)
			fmt.Println(i)
			x := utils.GetTotalLength(dataOut)
			fmt.Println("TotalPacket: ", x)
			tunFile.Write(dataOut)
		}()
	}

}
