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
		MaxIncomingStreams:    1000,
		MaxIncomingUniStreams: 1000,
		KeepAlivePeriod:       30,
		Versions:              []quic.VersionNumber{quic.Version2},
	}

	return quicConfig

}
