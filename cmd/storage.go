/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

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

type S3Object struct {
	Kind       string `yaml:"kind"`
	BucketName string `yaml:"name"`
}

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// yamlファイルに読み込み
		s3Object := S3Object{}
		b, _ := os.ReadFile("./sample/s3.yaml")
		yaml.Unmarshal(b, &s3Object)

		fmt.Printf("s3Object: %+v\n", s3Object)
		fmt.Printf("Tons: %v\n", s3Object.BucketName)

		// ローカルの認証情報を読み込む
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}

		// s3用のclientを作る
		client := s3.NewFromConfig(cfg)
		output, err := client.CreateBucket(
			context.TODO(),
			&s3.CreateBucketInput{
				Bucket: aws.String(s3Object.BucketName),
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
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
