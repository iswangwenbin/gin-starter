/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gin-starter",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GlobalConfig 全局配置变量
var GlobalConfig *configx.Config

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config/local.yaml", "config file (default is config/local.yaml)")
	rootCmd.PersistentFlags().StringP("env", "e", "local", "Set the environment.")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode")

	// 在每个命令执行前加载配置
	cobra.OnInitialize(initConfig)
}

// initConfig 加载配置文件
func initConfig() {
	// 获取环境参数
	env, _ := rootCmd.PersistentFlags().GetString("env")
	if env == "" {
		env = "local"
	}

	// 构建配置文件路径
	configPath := filepath.Join("config", env+".yaml")

	// 加载配置
	cfg, err := configx.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置全局配置
	GlobalConfig = cfg
}
