package repository

import (
	"context"

	"mailer/model"
)

type DBRepository interface {
	Scan(ctx context.Context) ([]model.User, error)
	Filter(users []model.User, filt model.Filter) []model.User
}
