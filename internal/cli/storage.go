/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sanoyo/aws-deployer/internal/session"

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
	cmd.Flags().BoolP("generate", "g", false, "generate yaml file")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	ok, err := cmd.Flags().GetBool("generate")
	if err != nil {
		return err
	}
	if ok {
		return generateStorageYaml(cmd)
	}

	// ファイルを読み込む
	filePath, err := cmd.Flags().GetString("yaml")
	if err != nil {
		return err
	}

	return createS3Bucket(ctx, filePath)
}

func createS3Bucket(ctx context.Context, filePath string) error {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	s3Ops := S3Ops{}
	err = yaml.Unmarshal(b, &s3Ops)
	if err != nil {
		return err
	}

	// 認証情報を取得
	session, err := session.NewProvider().Default()
	if err != nil {
		return err
	}
	client := s3.New(session)

	_, err = client.CreateBucketWithContext(
		ctx,
		&s3.CreateBucketInput{
			Bucket: aws.String(s3Ops.BucketName),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func generateStorageYaml(cmd *cobra.Command) error {

	var qs = []*survey.Question{
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

	s3Ops := S3Ops{
		Kind:       "S3",
		BucketName: answers.StorageName,
	}

	b, err := yaml.Marshal(s3Ops)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("./output/%s.yaml", answers.FileName), b, 0644)
	if err != nil {
		return err
	}

	return nil
}
