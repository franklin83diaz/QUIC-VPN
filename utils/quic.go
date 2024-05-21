package utils

import (
	"time"

	"github.com/quic-go/quic-go"
)

func GenerateQUICConfig() *quic.Config {

	// Configuración QUIC
	quicConfig := &quic.Config{
		MaxIdleTimeout:       60 * time.Second,
		HandshakeIdleTimeout: 60 * time.Second,
		KeepAlivePeriod:      30 * time.Second,
		// InitialStreamReceiveWindow: 20 * 1024 * 1024,
		// MaxConnectionReceiveWindow: 20 * 1024 * 1024,
	}

	return quicConfig

}
