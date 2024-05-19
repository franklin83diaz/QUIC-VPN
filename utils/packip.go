package utils

import (
	"encoding/binary"
	"errors"
)

func IsIPv4(packet []byte) bool {
	return packet[0]>>4 == 4
}

func GetTotalLength(packet []byte) uint16 {
	return binary.BigEndian.Uint16(packet[2:4])
}

func ValidateIPPacket(packet []byte) error {
	if len(packet) < 20 {
		return errors.New("packet is shorter than minimum IP header length")
	}

	ihl := packet[0] & 0x0F
	if ihl < 5 {
		return errors.New("length of IP header is less than 5 words")
	}

	headerLen := int(ihl * 4)
	if len(packet) < headerLen {

		return errors.New("packet is shorter than the length of the IP header")
	}

	totalLen := binary.BigEndian.Uint16(packet[2:4])
	if len(packet) < int(totalLen) {
		return errors.New("packet is shorter than the length of the IP packet")
	}

	if !validateChecksum(packet[:headerLen]) {
		return errors.New("invalid IP header checksum")
	}

	return nil
}

func validateChecksum(header []byte) bool {
	var sum uint32
	for i := 0; i < len(header); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(header[i:]))
	}

	for sum > 0xFFFF {
		sum = (sum >> 16) + (sum & 0xFFFF)
	}

	return ^uint16(sum) == 0
}
