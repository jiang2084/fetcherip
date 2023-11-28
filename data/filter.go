package data

import (
	"crypto/tls"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// CheckIp 验证代理是否有效
func CheckIp(ip *IP) bool {
	var testIp string
	var testUrl string

	if ip.ProxyType == "http" {
		testIp = fmt.Sprintf("http://%s:%d", ip.ProxyHost, ip.ProxyPort)
		testUrl = "http://httpbin.org/get"
	}
	// 解析代理地址
	proxy, parseErr := url.Parse(testIp)
	if parseErr != nil {
		logger.Errorf("parse error: %v\n", parseErr.Error())
		return false
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	// 设置网络传输
	netTransport := &http.Transport{
		DialContext: dialer.DialContext,
		Proxy:       http.ProxyURL(proxy),
		// true表示开启长连接
		MaxConnsPerHost:       20,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       20 * time.Second,
		ResponseHeaderTimeout: time.Second * time.Duration(10),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	_ = http2.ConfigureTransport(netTransport)
	// 创建客户端
	httpClient := &http.Client{
		Transport: netTransport,
	}
	begin := time.Now()
	// 使用代理ip访问地址
	res, err := httpClient.Get(testUrl)
	if err != nil {
		logger.Warnf("testIp: %s, testUrl: %s: error msg: %v", testIp, testUrl, err.Error())
		return false
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		speed := time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 // ms
		ip.ProxySpeed = int(speed)
		//UpdateIp(ip)
		return true
	}
	return false
}

// CheckDB 定期检查表里代理是否有效
func CheckDB() {
	nums := GetAllIPNums()
	logger.Infof("Before check, DB has: %d records.", nums)

	allIPs := GetAllIP()
	// 开了好多协程，这里可以优化
	var wg sync.WaitGroup
	for _, v := range allIPs {
		wg.Add(1)
		go func(ip IP) {
			if !CheckIp(&ip) {
				// 无法代理的ip直接删除
				DeleteIp(&ip)
			}
			wg.Done()
		}(v)
	}

	wg.Wait()
	nums = GetAllIPNums()
	logger.Infof("After check, DB has: %d records.", nums)
}

// CheckProxy 验证每个ip是否可用
func CheckProxy(ip *IP) {
	if CheckIp(ip) {
		logger.Infof("可用的IP:%s,端口:%v", ip.ProxyHost, ip.ProxyPort)
		SaveIp(ip)
	}
}
