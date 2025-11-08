package port

import (
	"context"
	"simple-blog-system/internal/app/user/model"
	"simple-blog-system/internal/app/user/payload"
)

type IUserService interface {
	Register(ctx context.Context, user model.AuthUserModel) (token string, err error)

	Login(ctx context.Context, user model.AuthUserModel) (token string, err error)

	GetUser(ctx context.Context, username string) (res *payload.User, err error)
}
