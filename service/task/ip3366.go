package task

import (
	"fetchip/data"
	"fetchip/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func Ip3366(ipChan chan<- *data.IP) []*data.IP {
	logger.Info("[Ip3366] fetch start")
	defer func() {
		recover()
		logger.Warnln("[Ip3366] fetch error")
	}()

	ips := make([]*data.IP, 0)

	for i := 20; i <= 40; i++ {
		url := fmt.Sprintf("http://www.ip3366.net/free/?stype=2&page=%v", i)
		fetchedPage, err := Get(url)
		if err != nil {
			// 本地ip被封了用代理ip请求
			logger.Info("正在使用代理抓取:")
			fetchedPage, err = ProxyGet(url)
			// 还是不行就过掉
			if err != nil {
				continue
			}
		}
		fetchedPage.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, s *goquery.Selection) {
				proxyIp := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
				proxyType := strings.TrimSpace(s.Find("td:nth-child(4)").Text())
				proxyLocation := strings.TrimSpace(s.Find("td:nth-child(5)").Text())
				fmt.Printf("ip:%s, host:%v, proxyType:%s\n", proxyIp, proxyPort, proxyType)
				ip := new(data.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = strings.ToLower(proxyType)
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = "http://www.ip3366.net/"
				ip.CreateTime = time.Now()
				ip.UpdateTime = time.Now()
				//ips = append(ips, ip)
				// 将数据放入管道
				ipChan <- ip
			})
		})
		sleepTime := utils.RandInt(0, 2)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	return ips
}
