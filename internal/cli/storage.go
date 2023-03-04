/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	internalAws "github.com/sanoyo/aws-deployer/internal/aws"
	"github.com/sanoyo/aws-deployer/internal/log"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

type s3Setting struct {
	Kind       string `yaml:"kind"`
	BucketName string `yaml:"name"`
}

type storageCommandOps struct {
	yaml         string
	generateFlag bool
}

type initStorageOpts struct {
	storageClient *internalAws.S3
	option        storageCommandOps
}

func BuildStorageCommand() *cobra.Command {
	ops := storageCommandOps{}
	stCmd := &cobra.Command{
		Use:   "storage",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		RunE: runCmdE(func(cmd *cobra.Command, args []string) error {
			opts, err := newInitStorageOpts(ops)
			if err != nil {
				return err
			}
			return run(opts)
		}),
	}

	stCmd.Flags().StringVar(&ops.yaml, "yaml", "", "Specify yaml file")
	stCmd.Flags().BoolVar(&ops.generateFlag, "generate", false, "generate yaml file")

	return stCmd
}

func newInitStorageOpts(ops storageCommandOps) (*initStorageOpts, error) {
	defaultSess, err := internalAws.NewProvider().Default()
	if err != nil {
		return nil, err
	}

	storage := initStorageOpts{
		storageClient: internalAws.NewStorage(defaultSess),
		option:        ops,
	}

	log.Logger.Info("successfully initialized")

	return &storage, nil
}

func (o *initStorageOpts) Validate() error {
	// どちらも指定されていない場合
	if o.option.yaml == "" && !o.option.generateFlag {
		return errors.New("please specify yaml file or generate flag")
	}

	// どちらも指定されている場合
	if o.option.yaml != "" && o.option.generateFlag {
		return errors.New("both flags are specified")
	}

	return nil
}

func (o *initStorageOpts) Execute() error {
	// generateオプションが指定されている場合
	if o.option.generateFlag {
		return o.generateStorageYaml()
	}

	// yamlファイルが指定されている場合
	return o.createS3Bucket()
}

func (o *initStorageOpts) createS3Bucket() error {
	log.Logger.Info("start to create s3")
	ctx := context.TODO()
	b, err := os.ReadFile(o.option.yaml)
	if err != nil {
		return err
	}

	ss := s3Setting{}
	err = yaml.Unmarshal(b, &ss)
	if err != nil {
		return err
	}

	_, err = o.storageClient.CreateBucketWithContext(
		ctx,
		&s3.CreateBucketInput{
			Bucket: aws.String(ss.BucketName),
		},
	)
	if err != nil {
		return err
	}

	log.Logger.Info("successfully to create s3")

	return nil
}

// FIXME: リファクタする.
func (o *initStorageOpts) generateStorageYaml() error {
	log.Logger.Info("start to generate s3 file")
	qs := []*survey.Question{
		{
			Name:      "storage_name",
			Prompt:    &survey.Input{Message: "Please input storage name?"},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
		{
			Name:      "file_name",
			Prompt:    &survey.Input{Message: "Please input output file name?"},
			Validate:  survey.Required,
			Transform: survey.ToLower,
		},
	}

	answers := struct {
		StorageName string `survey:"storage_name"`
		FileName    string `survey:"file_name"`
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	ss := s3Setting{
		Kind:       "S3",
		BucketName: answers.StorageName,
	}

	b, err := yaml.Marshal(ss)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("./output/%s.yaml", answers.FileName), b, 0o644)
	if err != nil {
		return err
	}

	log.Logger.Info("successfully to generate s3 file")

	return nil
}
