package nettest

import (
	"net"
	"strconv"
	"strings"
)

// GetTCP returns a TCP port that is available for use.
//
// NOTE: in the very rare case when no ports are available, this method will panic
func GetTCP() int {
	address := net.JoinHostPort("localhost", "0")
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	return getPortNo(listener.Addr().String())
}

// GetUDP returns a UDP port that is available for use.
//
// NOTE: in the very rare case when no ports are available, this method will panic
func GetUDP() int {
	address := net.JoinHostPort("localhost", "0")
	listener, err := net.ListenPacket("udp", address)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	return getPortNo(listener.LocalAddr().String())
}

func getPortNo(address string) int {
	chunks := strings.Split(address, ":")
	port, err := strconv.Atoi(chunks[len(chunks)-1])
	if err != nil {
		panic(err)
	}

	return port
}
