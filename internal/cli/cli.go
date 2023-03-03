package cli

import (
	"os"

	"github.com/spf13/cobra"
)

type cmd interface {
	// Validate returns an error if a flag's value is invalid.
	// Validate() error

	// Ask prompts for flag values that are required but not passed in.
	// Ask() error

	// Execute runs the command after collecting all required options.
	Execute() error
}

func runCmdE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 && args[0] == "help" {
			_ = cmd.Help() // Help always returns nil.
			os.Exit(0)
		}
		return f(cmd, args)
	}
}

func run(cmd cmd) error {
	// TODO: 後々どのcliも3つのメソッドすべてを実装するようにする予定
	// if err := cmd.Validate(); err != nil {
	// 	return err
	// }
	// if err := cmd.Ask(); err != nil {
	// 	return err
	// }
	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}
