package middleware

import (
	"fetchip/utils"
	"github.com/fsnotify/fsnotify"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var (
	Vip           = viper.New()
	ConfigFile    = ""
	ServerSetting = new(YamlSetting)
)

/*
redis:
    host: 127.0.0.1
    port: 6379
*/

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

/*
	host: 127.0.0.1
    port: 3306
 	dbName: ips
    username: root
    password: 123456
    charset: utf8mb4
    maxIdleConns: 5
    maxOpenConns: 100
*/

type Database struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DbName       string `yaml:"dbName"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Charset      string `yaml:"charset"`
	MaxIdleConns string `yaml:"maxIdleConns"`
	MaxOpenConns string `yaml:"maxOpenConns"`
}

/*
	filePath:logs
    fileName: ip.log
    # 可以是控制台或者文件，默认是控制台，多模式使用','分割
    mode:file
*/

type Logs struct {
	FilePath string `yaml:"filePath"`
	FileName string `yaml:"fileName"`
	Mode     string `yaml:"mode"`
}

type YamlSetting struct {
	Database Database `yaml:"mysql" mapstructure:"mysql"`
	Redis    Redis    `yaml:"redis"`
	Logs     Logs     `yaml:"logs"`
}

func InitConfig() {
	if ConfigFile != "" {
		if !utils.PathExists(ConfigFile) {
			logger.Errorf("没有这个文件或目录: %s", ConfigFile)
			os.Exit(-1)
		} else {
			Vip.SetConfigFile(ConfigFile)
			Vip.SetConfigType("yaml")
		}
	} else {
		logger.Errorf("配置文件不能为空: %s", ConfigFile)
		os.Exit(-1)
	}

	err := Vip.ReadInConfig()
	if err != nil {
		logger.Errorf("无法读取配置文件:%s", err)
	}

	Vip.WatchConfig()
	Vip.OnConfigChange(func(in fsnotify.Event) {
		logger.Infof("配置文件内容发生改变:%s", ConfigFile)
		ServerSetting = GetConfig(Vip)
	})
	Vip.AllSettings()
	ServerSetting = GetConfig(Vip)
}

func GetConfig(vip *viper.Viper) *YamlSetting {
	setting := new(YamlSetting)
	// 解析配置文件，反序列化
	if err := vip.Unmarshal(setting); err != nil {
		logger.Errorf("解析失败:%s", err)
		os.Exit(-1)
	}
	return setting
}
