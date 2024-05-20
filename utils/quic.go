package utils

import (
	"time"

	"github.com/quic-go/quic-go"
)

func GenerateQUICConfig() *quic.Config {

	// Configuración QUIC
	quicConfig := &quic.Config{
		MaxIdleTimeout:        60 * time.Second,
		HandshakeIdleTimeout:  60 * time.Second,
		MaxIncomingStreams:    1,
		MaxIncomingUniStreams: 1,
		KeepAlivePeriod:       30,
		// InitialStreamReceiveWindow: 20 * 1024 * 1024,
		// MaxConnectionReceiveWindow: 20 * 1024 * 1024,
	}

	return quicConfig

}
