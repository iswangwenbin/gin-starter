/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/iswangwenbin/gin-starter/internel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取环境参数
		debug, _ := cmd.Flags().GetBool("debug")
		env := viper.GetString("env")
		if env == "" {
			env = "development"
		}
		options := internal.WithDefaults()
		if debug {
			options = internal.WithDebug()
		}
		
		server, err := internal.NewServer(env, options...)
		if err != nil {
			log.Fatalf("Failed to create server: %v", err)
		}

		fmt.Printf("Starting server in %s mode...\n", env)
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
