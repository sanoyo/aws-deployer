package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Provider struct {
	defaultSess *session.Session
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Default() (*session.Session, error) {
	if p.defaultSess != nil {
		return p.defaultSess, nil
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
		Profile: "default",
	})
	if err != nil {
		return nil, err
	}

	p.defaultSess = sess
	return p.defaultSess, nil
}
