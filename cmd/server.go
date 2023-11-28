/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fetchip/middleware"
	"fmt"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:          "start",
	SilenceUsage: true,
	Short:        "Get Application config info",
	Example:      "fetchip config -f conf/config.yml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
	},
}

func init() {
	cobra.OnInitialize(middleware.InitConfig)

	//setting := middleware.ServerSetting

	serverCmd.PersistentFlags().StringVarP(&middleware.ConfigFile, "configFile", "c", "conf/config.yml", "config file")

	// 必须配置
	_ = serverCmd.MarkFlagRequired("configFile")

	rootCmd.AddCommand(serverCmd)

}
