package task

import (
	"crypto/tls"
	"errors"
	"fetchip/data"
	"fetchip/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

func ProxyGet(urls string) (*goquery.Document, error) {
	var doc *goquery.Document
	ip := data.GetOneIp()
	logger.Infof("使用的代理ip是:%s, %v", ip.ProxyHost, ip.ProxyPort)
	if ip.ProxyId == 0 {
		return doc, errors.New("代理ip池里面没有可用ip")
	}
	var testIp string
	var testUrl string

	if ip.ProxyType == "http" {
		testIp = fmt.Sprintf("http://%s:%d", ip.ProxyHost, ip.ProxyPort)
		testUrl = urls
	}
	// 解析代理地址
	proxy, parseErr := url.Parse(testIp)
	if parseErr != nil {
		logger.Errorf("parse error: %v\n", parseErr.Error())
		return doc, errors.New("解析代理地址出错")
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
	// 使用代理ip访问地址
	resp, err := httpClient.Get(testUrl)
	if err != nil {
		logger.Warnf("testIp: %s, testUrl: %s: error msg: %v", testIp, testUrl, err.Error())
		return doc, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			logrus.Errorf("charset convert failed: %v", err)
			return doc, err
		}

		doc, err = goquery.NewDocumentFromReader(content)
		if err != nil {
			logrus.Errorf("goquery http response body reader error: %v", err)
			return doc, err
		}
	}
	return doc, nil
}

func Get(url string) (*goquery.Document, error) {
	logrus.Infof("Fetch url: %s", url)
	var doc *goquery.Document
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return doc, err
	}
	client := &http.Client{
		Jar:     cookieJar,
		Timeout: 10 * time.Second,
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Proxy-Switch-Ip", "yes")
	req.Header.Set("User-Agent", utils.GetUserAgent())
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.")
	req.Header.Set("Accept-language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		return doc, err
	}

	if resp.StatusCode == http.StatusOK {
		content, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if err != nil {
			logrus.Errorf("charset convert failed: %v", err)
			return doc, err
		}

		doc, err = goquery.NewDocumentFromReader(content)
		if err != nil {
			logrus.Errorf("goquery http response body reader error: %v", err)
			return doc, err
		}
	}

	return doc, nil
}
