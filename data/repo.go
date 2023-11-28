package data

import (
	"fetchip/middleware"
	"fetchip/utils"
	logger "github.com/sirupsen/logrus"
)

// SaveIp 保存ip
func SaveIp(ip *IP) {
	db := middleware.GetDB().Begin()
	ipModel := GetIpByHost(ip.ProxyHost)
	if ipModel.ProxyHost == "" {
		// 说明表里没有数据
		result := db.Create(ip)
		if result.Error != nil {
			logger.Errorf("save ip: %s, error msg: %v", ip.ProxyHost, result.Error)
			db.Rollback()
		}
	} else {
		UpdateIp(ipModel)
	}
	db.Commit()
}

func UpdateIp(ip *IP) {
	db := middleware.GetDB().Begin()
	ipMap := make(map[string]interface{}, 0)
	ipMap["proxy_speed"] = ip.ProxySpeed
	ipMap["update_time"] = utils.FormatDateTime()
	if ip.ProxyId != 0 {
		result := db.Model(new(IP)).Where("proxy_id = ?", ip.ProxyId).Updates(ipMap)
		if result.Error != nil {
			logger.Errorf("update ip: %s, error msg: %v", ip.ProxyHost, result.Error)
			db.Rollback()
		}
	}
	db.Commit()
}

func DeleteIp(ip *IP) {
	db := middleware.GetDB().Begin()
	result := db.Model(new(IP)).Where("proxy_id = ?", ip.ProxyId).Delete(ip)
	if result.Error != nil {
		logger.Errorf("delete ip: %s, error msg: %v", ip.ProxyHost, result.Error)
		db.Rollback()
	}
	db.Commit()
}

// GetIpByHost 根据host获取ip信息
func GetIpByHost(host string) *IP {
	db := middleware.GetDB()
	ipModel := new(IP)
	db.Model(&ipModel).Where("proxy_host = ?", host).Find(ipModel)
	return ipModel
}

// GetAllIPNums 查询表里数据
func GetAllIPNums() int64 {
	db := middleware.GetDB()
	var count int64
	result := db.Model(new(IP)).Count(&count)
	if result.Error != nil {
		logger.Errorf("ip count: %d, error msg:%v", count, result.Error)
		return -1
	}
	return count
}

// GetAllIP 获取所有ip数据
func GetAllIP() []IP {
	db := middleware.GetDB()
	list := make([]IP, 0)
	result := db.Model(new(IP)).Find(&list)
	ipCount := len(list)
	if result.Error != nil {
		logger.Warnf("ip count: %d, error msg: %v\n", ipCount, result.Error)
		return nil
	}
	return list
}

// GetOneIp 随机获取一个ip
func GetOneIp() *IP {
	oneIp := new(IP)
	lists := GetAllIP()
	size := len(lists)
	if len(lists) != 0 {
		idx := utils.RandInt(0, size)
		oneIp = &lists[idx]
	}
	return oneIp
}
