package task

import (
	"fetchip/middleware"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	//url := "https://www.cnblogs.com"
	//doc := Get(url)
	//doc.Find("")

	//fmt.Println(doc)
}

func TestProxyGet(t *testing.T) {
	path := "/Users/changba-os/Desktop/github/go/mine/fetchip/conf/config.yml"

	config := middleware.InitConfig(path)

	middleware.InitDB(&config.Database)

	//_, err := ProxyGet("http://www.ip3366.net/free/?stype=1&page=2")
	//fmt.Println(err)
	_, err := Get("http://www.ip3366.net/free/?stype=1&page=2")
	fmt.Println(err)
}
