package task

import (
	data2 "fetchip/data"
	"fmt"
	"sync"
	"testing"
)

func TestIp89(t *testing.T) {
	ips := Ip89()
	for _, ip := range ips {
		//fmt.Println(ip.ProxyHost)
		go func(ip *data2.IP) {
			fmt.Println(ip.ProxyHost)
			res := data2.CheckIp(ip)
			if res {
				fmt.Printf("ip:%s, port:%v", ip.ProxyHost, ip.ProxyPort)
			}
		}(ip)
	}
	for {

	}
}

func TestIp3366(t *testing.T) {
	ips := Ip3366()
	var wg sync.WaitGroup
	for _, ip := range ips {
		wg.Add(1)
		go func(ip *data2.IP) {
			res := data2.CheckIp(ip)
			fmt.Printf("ip:%s, host:%v res: %v\n", ip.ProxyHost, ip.ProxyPort, res)
			if res {
				fmt.Printf("ip:%s, port:%v", ip.ProxyHost, ip.ProxyPort)
			}
			wg.Done()
		}(ip)
	}
	wg.Wait()
}
