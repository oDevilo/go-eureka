package eureka_client

import (
	"net"
	"strings"
)

// 获取本地ip
func GetLocalIP() (string, error) {
	addressArr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addressArr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

// url拼接
func UrlAppend(start, end string) string {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)
	return strings.TrimRight(start, "/") + "/" + strings.TrimLeft(end, "/")
}
