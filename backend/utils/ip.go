package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

// GetIpToAddress 根据ip地址获取到地址
func GetIpToAddress(ip string) *CityInfo {
	var cityInfo = &CityInfo{}
	var ipStr string
	var infos []string
	p, _ := NewIpdb("./config/qqzeng-ip-utf8.dat")
	ipStr = p.Get(ip)
	infos = strings.Split(ipStr, "|")
	if infos[1] == "保留" {
		cityInfo.Country = "Internal network"
		return cityInfo
	}
	cityInfo.Country = infos[1]
	cityInfo.Province = infos[2]
	cityInfo.City = infos[3]
	return cityInfo
}

// GetCurrentIP 获取当前请求ip
func GetCurrentIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = ctx.Request.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func GetServiceIpLocal() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("无法获取本地IP地址")
}

func GetServiceIp() (string, error) {
	return GetServiceIpLocal()
}
