/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gin-starter",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringP("config", "c", "config/local.yaml", "config file (default is config/local.yaml)")
	rootCmd.PersistentFlags().StringP("env", "e", "", "Set the environment.")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode")

	cobra.OnInitialize(func() {
		// 获取命令行参数
		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		envValue, _ := rootCmd.PersistentFlags().GetString("env")
		configDir := filepath.Dir(configPath)
		initConfig(envValue, configDir)
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig(env, path string) {
	var name string
	// 优先使用环境变量 APP_ENV
	confDir := os.Getenv("APP_ENV")
	if confDir != "" {
		path = confDir
	}
	if path == "" {
		exePath, _ := os.Executable()
		path = filepath.Dir(exePath) + "/config"
	}
	// 根据环境选择配置文件名
	switch env {
	case "production":
		name = "production"
	case "development":
		name = "development"
	default:
		if _, err := os.Stat(path + "/local.yaml"); err == nil {
			name = "local"
		} else {
			name = "development"
		}
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config file error: %v", err)
	} else {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	// 配置热重载及回调
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}
