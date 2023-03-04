/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/sanoyo/aws-deployer/internal/cli"
	"github.com/sanoyo/aws-deployer/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "aws-deployer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var Logger *zap.Logger

func init() {
	cobra.OnInitialize(
		initLogging,
	)
	rootCmd.AddCommand(cli.BuildStorageCommand())
}

func main() {
	Execute()
}

func initLogging() {
	log.NewLogger()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
