package data

import (
	"fetchip/middleware"
	"fmt"
	"testing"
	"time"
)

// 两个问题，一个是编码，一个事务，注意执行的返回值，Error是包在里面的
func TestSave(t *testing.T) {
	config := middleware.ServerSetting

	db := middleware.InitDB(&config.Database)

	// 2个问题，表名和数据的问题
	db.Set("gorm:table_options", "CHARSET=utf8")
	err := db.Migrator().CreateTable(&IP{})
	if err != nil {
		return
	}

	//user := new(User)
	//user.Name = "jiang"
	//db.Create(user)

	//ip := new(IP)
	//ip.ProxyHost = "127.0.0.2"
	//ip.ProxyPort = 3306
	//ip.ProxyType = "http"
	//ip.ProxyLocation = "测试2"
	//ip.ProxySpeed = 100
	//ip.ProxySource = "http://www.66ip.cn"
	//ip.CreateTime = time.Now()
	//ip.UpdateTime = time.Now()
	//
	//db.Create(ip)
}

func TestSaveIp(t *testing.T) {
	config := middleware.ServerSetting
	middleware.InitDB(&config.Database)
	ip := new(IP)
	ip.ProxyHost = "127.0.0.1"
	ip.ProxyPort = 3306
	ip.ProxyType = "http"
	ip.ProxyLocation = "本地weee"
	ip.ProxySpeed = 100
	ip.ProxySource = "http://www.66ip.cn"
	ip.CreateTime = time.Now()
	ip.UpdateTime = time.Now()

	SaveIp(ip)
}

func TestGetOneIp(t *testing.T) {

	config := middleware.ServerSetting
	middleware.InitDB(&config.Database)

	ip := GetOneIp()
	fmt.Println(ip.ProxyHost, ip.ProxyPort)
}
