package port

import (
	"context"
	"simple-blog-system/internal/app/user/model"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, user model.AuthUserModel) (model.AuthUserModel, error)

	GetUserByUsername(ctx context.Context, username string) (user []model.AuthUserModel, err error)

	GetPasswordByUsername(ctx context.Context, username string) (user []model.AuthUserModel, err error)

	UpdateLastLogin(ctx context.Context, user model.AuthUserModel) error
}
