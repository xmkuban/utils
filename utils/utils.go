package utils

import (
	"crypto/md5"
	"fmt"
	"io"
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

func MD5(src string) string {
	h := md5.New()
	io.WriteString(h, src)
	return fmt.Sprintf("%x", h.Sum(nil))
}
func MD5ByLength16(src string) string {
	res := MD5(src)
	res = res[8:]
	res = res[:16]
	return res
}

func EmptyString(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}
