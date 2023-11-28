package service

import (
	data2 "fetchip/data"
	"fetchip/service/task"
	logger "github.com/sirupsen/logrus"
	"sync"
	"time"
)

func Task() {
	// 检验已经有的数据
	ipChan := make(chan *data2.IP, 2000)

	go func() {
		data2.CheckDB()
	}()

	// 开启50个协程验证抓取下来的ip是否可用
	for i := 0; i < 50; i++ {
		go func() {
			for {
				data2.CheckProxy(<-ipChan)
			}
		}()
	}

	for {
		nums := data2.GetAllIPNums()
		logger.Printf("Chan: %v, IP: %d\n", len(ipChan), nums)
		// 如果IP池数量小于100就抓取
		if len(ipChan) < 2000 {
			go Do(ipChan)
		}
		time.Sleep(300 * time.Second)
	}
}

func Do(ipChan chan<- *data2.IP) {
	var wg sync.WaitGroup
	siteFuncList := []func(ipChan chan<- *data2.IP) []*data2.IP{
		task.Ip89,
		//task.Ip3366,
	}

	for _, site := range siteFuncList {
		wg.Add(1)
		go func(site func(ipChan chan<- *data2.IP) []*data2.IP) {
			site(ipChan)
			//obtainIPs := site(ipChan)
			//logger.Infof("obtainIPs:", len(obtainIPs))
			//for _, ip := range obtainIPs {
			//	ipChan <- ip
			//}
		}(site)
	}
}
