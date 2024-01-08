package service

import (
	"context"

	"mailer/model"
	"mailer/repository"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type sesClient struct {
	client *sesv2.Client
}

func NewSESClient(ctx context.Context) (repository.SESRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, err
	}
	client := sesv2.NewFromConfig(cfg)
	return &sesClient{client: client}, nil
}

func (c *sesClient) SendEmail(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error) {
	return c.client.SendEmail(ctx, input)
}

func (c *sesClient) AssembleEmail(messege *model.MailMessege) *sesv2.SendEmailInput {
	return &sesv2.SendEmailInput{
		FromEmailAddress: &messege.Sender,
		Destination: &types.Destination{
			ToAddresses: []string{messege.Recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data: &messege.Body,
					},
				},
				Subject: &types.Content{
					Data: &messege.Subject,
				},
			},
		},
	}
}
