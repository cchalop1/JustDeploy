package adapter

import (
	"fmt"
	"net"
)

type NetworkAdapter struct {
}

func NewNetworkAdapter() *NetworkAdapter {
	return &NetworkAdapter{}
}

func (n *NetworkAdapter) GetServerURL(port string) (string, error) {
	ip, err := n.getCurrentIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%s", ip, port), nil
}

func (n *NetworkAdapter) getCurrentIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", nil
}
