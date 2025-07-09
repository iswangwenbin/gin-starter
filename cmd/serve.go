/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/iswangwenbin/gin-starter/internal"
	"github.com/iswangwenbin/gin-starter/pkg/configx"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long: `Start the server with the specified configuration.

Examples:
  gin-starter serve                    # Start with default settings
  gin-starter serve --env production   # Start in production mode
  gin-starter serve --env local        # Start with local configuration
  gin-starter serve --debug            # Start with debug enabled`,

	Run: func(cmd *cobra.Command, args []string) {
		// 获取环境参数
		debug, _ := cmd.Flags().GetBool("debug")
		env, _ := cmd.Flags().GetString("env")
		if env == "" {
			env = "development"
		}

		// 加载配置文件
		configPath := filepath.Join("config", env+".yaml")
		cfg, err := configx.Load(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// 选择服务器选项
		options := internal.WithDefaults()
		if debug {
			options = internal.WithDebug()
		}

		server, err := internal.NewServer(env, options...)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		fmt.Printf("Starting server in %s mode on %s...\n", env, cfg.GetServerAddress())
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("env", "development", "Environment (development, production, local)")
	serveCmd.Flags().Bool("debug", false, "Enable debug mode")
}
