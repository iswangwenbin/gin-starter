package cmd

import (
	"fmt"
	"log"

	"github.com/iswangwenbin/gin-starter/internal/core"
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
		// 获取环境参数（使用全局标志）
		debug, _ := cmd.Parent().PersistentFlags().GetBool("debug")
		env, _ := cmd.Parent().PersistentFlags().GetString("env")

		// 使用全局配置（已在 root.go 中加载）
		cfg := GlobalConfig
		if cfg == nil {
			log.Fatalf("Global config not loaded")
		}

		// 选择服务器选项
		options := core.WithDefaults()
		if debug {
			options = core.WithDebug()
		}

		// 创建服务器
		server, err := core.NewServer(env, options...)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		// 创建生命周期管理器并运行
		lifecycle := core.NewLifecycle(server)
		fmt.Printf("Starting server in %s mode on %s...\n", env, cfg.GetServerAddress())

		if err := lifecycle.Run(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
