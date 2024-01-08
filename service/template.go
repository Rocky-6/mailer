package service

import (
	"context"
	"encoding/json"
	"io"

	"mailer/model"
	"mailer/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3Client struct {
	client S3API
}

func NewS3Client(ctx context.Context) (repository.S3Repository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, err
	}
	return &s3Client{client: s3.NewFromConfig(cfg)}, nil
}

func (c *s3Client) FetchTemplate(ctx context.Context, bucket string, filename string) (*model.Template, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	buff, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	template := &model.Template{}
	err = json.Unmarshal(buff, template)
	if err != nil {
		return nil, err
	}

	return template, nil
}
