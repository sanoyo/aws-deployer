/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	internalAws "github.com/sanoyo/aws-deployer/internal/aws"
	"github.com/sanoyo/aws-deployer/internal/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type taskSetting struct {
	Cluster    string  `yaml:"cluster"`
	Name       string  `yaml:"name"`
	LaunchType string  `yaml:"launchType"`
	Network    Network `yaml:"network"`
}

type Network struct {
	Subnets []string `yaml:"subnets"`
}

type taskCommandOps struct {
	yaml         string
	generateFlag bool
}

type initTaskOpts struct {
	taskClient *internalAws.Task
	option     taskCommandOps
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
	defaultSess, err := internalAws.NewProvider().Default()
	if err != nil {
		return nil, err
	}

	task := initTaskOpts{
		taskClient: internalAws.NewTask(defaultSess),
		option:     ops,
	}

	log.Logger.Info("successfully initialized")

	return &task, nil
}

func (o *initTaskOpts) Validate() error {
	// どちらも指定されていない場合
	if o.option.yaml == "" && !o.option.generateFlag {
		return errors.New("please specify yaml file or generate flag")
	}

	// どちらも指定されている場合
	if o.option.yaml != "" && o.option.generateFlag {
		return errors.New("both flags are specified")
	}

	log.Logger.Info("successfully to validate ecs task")
	return nil
}

func (o *initTaskOpts) Execute() error {
	b, err := os.ReadFile(o.option.yaml)
	if err != nil {
		return err
	}

	task := taskSetting{}
	err = yaml.Unmarshal(b, &task)
	if err != nil {
		return err
	}

	_, err = o.taskClient.RunTask(
		&ecs.RunTaskInput{
			Cluster:        aws.String(task.Cluster),
			TaskDefinition: aws.String(task.Name),
			LaunchType:     aws.String(o.whereLaunchType(task.LaunchType)),
			NetworkConfiguration: &ecs.NetworkConfiguration{
				AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
					Subnets: aws.StringSlice(task.Network.Subnets),
				},
			},
		})
	if err != nil {
		return err
	}

	log.Logger.Info("successfully to create ecs task")
	return nil
}

func (o *initTaskOpts) whereLaunchType(lanchType string) string {
	switch lanchType {
	case "fargate":
		return ecs.LaunchTypeFargate
	case "ec2":
		return ecs.LaunchTypeEc2
	case "external":
		return ecs.LaunchTypeExternal
	default:
		return ""
	}
}
