/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iswangwenbin/gin-starter/internal/core"
	"github.com/iswangwenbin/gin-starter/internal/worker"
	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the install event worker",
	Long: `Start the install event worker to process events from Redis Stream.

The worker consumes install events from Redis Stream and writes them to ClickHouse.
It runs independently from the main server process.

Examples:
  gin-starter worker                   # Start with default settings
  gin-starter worker --env production  # Start in production mode
  gin-starter worker --debug           # Start with debug enabled`,

	Run: func(cmd *cobra.Command, args []string) {
		// 获取环境参数（使用全局标志）
		debug, _ := cmd.Parent().PersistentFlags().GetBool("debug")
		env, _ := cmd.Parent().PersistentFlags().GetString("env")

		// 使用全局配置（已在 root.go 中加载）
		cfg := GlobalConfig
		if cfg == nil {
			log.Fatalf("Global config not loaded")
		}

		// 创建核心服务器（只启用必要的服务）
		options := []core.Option{core.StartCache, core.StartClickHouse}
		if debug {
			options = append(options, core.StartDebug)
		}

		server, err := core.NewServer(env, options...)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		// 初始化服务器依赖（Redis、ClickHouse）
		if err := server.InitDependencies(); err != nil {
			log.Fatalf("Failed to initialize server dependencies: %v", err)
		}

		// 创建 Worker
		installEventWorker := worker.NewInstallEventWorker(
			server.Cache,
			server.ClickHouse,
			server.Logger(),
		)

		fmt.Printf("Starting install event worker in %s mode...\n", env)

		// 启动 Worker
		if err := installEventWorker.Start(); err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
		defer installEventWorker.Stop()

		// 等待停止信号
		server.Logger().Info("Install event worker started, waiting for signals...")
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		server.Logger().Info("Received shutdown signal, stopping worker...")
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}