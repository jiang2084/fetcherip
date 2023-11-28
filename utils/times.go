package utils

import "time"

// FormatDateTime 格式化日期
func FormatDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
