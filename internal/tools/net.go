package tools

import (
	"fmt"
	"net"
)

func TrustedSubnet(ipStr string, subnet net.IPNet) (bool, error) {

	// смотрим заголовок запроса X-Real-IP
	// парсим ip
	ip := net.ParseIP(ipStr)

	if ip == nil {
		return false, fmt.Errorf("failed parse ip from http header %s", ipStr)
	}
	if !subnet.Contains(ip) {
		return false, fmt.Errorf("IP doesn't trusted %s", ip.String())
	}
	return true, nil

}
