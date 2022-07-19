package auth

import (
	"context"

	"dclibrary.com/internals/app/models"
)

type Repository interface {
	Insert(ctx context.Context, user *models.UserAuth) error
	Get(ctx context.Context, username, password string) (*models.User, error)
}
