package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3API interface {
	CreateBucketWithContext(ctx context.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error)
}

type S3 struct {
	s3Client s3API
}

func NewStorage(s *session.Session) *S3 {
	return &S3{
		s3Client: s3.New(s),
	}
}

func (s *S3) CreateBucketWithContext(ctx context.Context, input *s3.CreateBucketInput, opts ...request.Option) (*s3.CreateBucketOutput, error) {
	return s.s3Client.CreateBucketWithContext(ctx, input, opts...)
}
