package repository

import (
	"context"

	"mailer/model"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type SESRepository interface {
	AssembleEmail(messege *model.MailMessege) *sesv2.SendEmailInput
	SendEmail(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error)
}
