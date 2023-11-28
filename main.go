/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fetchip/cmd"
	"fetchip/middleware"
	"fetchip/service"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 检查命令行参数
	cmd.Execute()

	setting := middleware.ServerSetting

	middleware.InitLog(&setting.Logs)

	middleware.InitDB(&setting.Database)

	service.Task()
}
