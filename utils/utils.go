package utils

import (
	"net"
	"strings"

	"github.com/google/uuid"
)

//GetDockerLocalIP  获取本地ip
func GetDockerLocalIP() string {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return ""
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		// handle err
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				return v.IP.String()
			case *net.IPAddr:
				return v.IP.String()
			}
		}
	}
	return ""
}

func GetUUID() string {
	return uuid.New().String()
}

func GetUUIDNoDash() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
