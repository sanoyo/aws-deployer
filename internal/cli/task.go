/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	internalAws "github.com/sanoyo/aws-deployer/internal/aws"
	"github.com/sanoyo/aws-deployer/internal/log"
	"github.com/spf13/cobra"
)

type taskSetting struct {
	Kind       string `yaml:"kind"`
	BucketName string `yaml:"name"`
}

type taskCommandOps struct {
	yaml         string
	generateFlag bool
}

type initTaskOpts struct {
	storageClient *internalAws.S3
	option        taskCommandOps
}

func BuildTaskCommand() *cobra.Command {
	ops := taskCommandOps{}
	var taskCmd = &cobra.Command{
		Use:   "task",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newInitTaskOpts(ops)
			if err != nil {
				return err
			}
			return run(opts)
		}),
	}

	taskCmd.Flags().StringVar(&ops.yaml, "yaml", "", "Specify yaml file")
	taskCmd.Flags().BoolVar(&ops.generateFlag, "generate", false, "generate yaml file")

	return taskCmd
}

func newInitTaskOpts(ops taskCommandOps) (*initTaskOpts, error) {
	return nil, nil
}

func (o *initTaskOpts) Validate() error {
	log.Logger.Info("successfully to validate ecs task")
	return nil
}

func (o *initTaskOpts) Execute() error {
	log.Logger.Info("successfully to create ecs task")
	return nil
}
