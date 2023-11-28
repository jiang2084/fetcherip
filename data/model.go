package data

import (
	"time"
)

// IP struct
type IP struct {
	ProxyId       int64     `gorm:"primary_key; auto_increment; not null" json:"id"`
	ProxyHost     string    `gorm:"type:varchar(255); not null; unique" json:"proxyHost"`
	ProxyPort     int       `gorm:"type:int(11); not null;" json:"proxyPort"`
	ProxyType     string    `gorm:"type:varchar(64); not null" json:"proxyType"`
	ProxyLocation string    `gorm:"type:varchar(255); default null" json:"proxyLocation"`
	ProxySpeed    int       `gorm:"type:int(20); not null; default 0" json:"proxySpeed"`
	ProxySource   string    `gorm:"type:varchar(64); not null;" json:"proxySource"`
	CreateTime    time.Time `gorm:"type:datetime;" json:"createTime"`
	UpdateTime    time.Time `gorm:"type:datetime;" json:"updateTime"`
}
