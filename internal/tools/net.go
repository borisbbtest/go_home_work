package tools

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func TrustedSubnet(r *http.Request, sub string) (bool, error) {

	// смотрим заголовок запроса X-Real-IP
	ipStr := r.Header.Get("X-Real-IP")
	_, subnet, _ := net.ParseCIDR(sub)
	// парсим ip
	ip := net.ParseIP(ipStr)
	if ip == nil {
		// если заголовок X-Real-IP пуст, пробуем X-Forwarded-For
		// этот заголовок содержит адреса отправителя и промежуточных прокси
		// в виде 203.0.113.195, 70.41.3.18, 150.172.238.178
		ips := r.Header.Get("X-Forwarded-For")
		fmt.Println(ips)
		// разделяем цепочку адресов
		ipStrs := strings.Split(ips, ",")
		// интересует только первый
		ipStr = ipStrs[0]
		// парсим
		ip = net.ParseIP(ipStr)
	}
	if ip == nil {
		return false, fmt.Errorf("failed parse ip from http header %s", ipStr)
	}
	if !subnet.Contains(ip) {
		return false, fmt.Errorf("IP doesn't trusted", ip.String())
	}
	return true, nil

}
