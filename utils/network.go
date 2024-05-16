package utils

import (
	"fmt"
	"net"
	"strings"
)

// Function to calculate available IPs from a CIDR mask
func AvailableIPs(cidr string) (int, string, []string, error) {
	// Split the CIDR into IP and prefix parts
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, "", nil, err
	}

	// Get the prefix length from the IP network
	prefixLength, _ := ipNet.Mask.Size()

	// Calculate the number of available IPs
	totalIPs := (1 << (32 - prefixLength))
	availableIPs := totalIPs - 3 // Subtract 3 server, network and broadcast addresses

	// Ensure the IP is an IPv4 address
	if strings.Contains(ip.String(), ":") {
		return 0, "", nil, fmt.Errorf("IPv6 addresses are not supported")
	}

	// Generate the list of available IPs
	var ipList []string
	for i := 1; i < totalIPs-1; i++ {
		ip := make(net.IP, len(ipNet.IP))
		copy(ip, ipNet.IP)
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j] += byte(i >> (8 * (len(ip) - 1 - j)))
		}
		ipList = append(ipList, ip.String())
	}
	ipServer := ipList[0]

	//remove first ip
	ipList = ipList[1:]

	return availableIPs, ipServer, ipList, nil
}
