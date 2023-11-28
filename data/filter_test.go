package data

import "testing"

func TestCheckIp(t *testing.T) {
	ip := new(IP)

	ip.ProxyType = "http"
	//ip.ProxyHost = "127.0.0.11"
	ip.ProxyHost = "58.20.77.156"
	ip.ProxyPort = 2323
	//ip.ProxyPort = 7890

	CheckIp(ip)
}
