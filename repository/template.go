package repository

import (
	"context"

	"mailer/model"
)

type S3Repository interface {
	FetchTemplate(ctx context.Context, bucket string, filename string) (*model.Template, error)
}
