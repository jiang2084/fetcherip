package task

import (
	"fetchip/data"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func Ip89(ipChan chan<- *data.IP) []*data.IP {
	logger.Info("[89ip] fetch start")
	defer func() {
		recover()
		logger.Warnln("[89ip] fetch error")
	}()

	ips := make([]*data.IP, 0)

	indexUrl := "https://www.89ip.cn/"
	for i := 1; i < 66; i++ {
		url := fmt.Sprintf("%s/index_%v.html", indexUrl, i)
		fetchedPage, err := Get(url)
		if err != nil {
			// 本地ip被封了用代理ip请求
			fetchedPage, err = ProxyGet(url)
			fmt.Println(err)
			// 还是不行就过掉
			if err != nil {
				continue
			}
		}
		fetchedPage.Find("table > tbody").Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(i int, s *goquery.Selection) {
				proxyIp := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
				proxyPort := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
				proxyLocation := strings.TrimSpace(s.Find("td:nth-child(3)").Text())

				ip := new(data.IP)
				ip.ProxyHost = proxyIp
				ip.ProxyPort, _ = strconv.Atoi(proxyPort)
				ip.ProxyType = "http"
				ip.ProxyLocation = proxyLocation
				ip.ProxySpeed = 100
				ip.ProxySource = indexUrl
				ip.CreateTime = time.Now()
				ip.UpdateTime = time.Now()
				//ips = append(ips, ip)

				ipChan <- ip
			})
		})
	}
	return ips
}
