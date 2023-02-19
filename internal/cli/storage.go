/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type S3Ops struct {
	Kind       string `yaml:"kind"`
	BucketName string `yaml:"name"`
}

func BuildStorageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
		and usage of using your command. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		RunE: run,
	}

	cmd.Flags().StringP("yaml", "y", "", "Specify yaml file")
	cmd.MarkFlagRequired("yaml")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	// ファイルを読み込む
	filePath, err := cmd.Flags().GetString("yaml")
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	s3Ops := S3Ops{}
	err = yaml.Unmarshal(b, &s3Ops)
	if err != nil {
		return err
	}

	// yamlファイルに読み込み
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)

	output, err := client.CreateBucket(
		ctx,
		&s3.CreateBucketInput{
			Bucket: aws.String(s3Ops.BucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(cfg.Region),
			},
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("output", output)

	return nil
}
