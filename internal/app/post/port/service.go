package port

import (
	"context"
	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/payload"
)

type IPostService interface {
	AddPost(ctx context.Context, username string, param payload.PostRequest) (res *model.PostModel, err error)
	UpdatePost(ctx context.Context, username string, id string, param payload.PostRequest) (res *model.PostModel, err error)
	DeletePost(ctx context.Context, username string, id string) (res *model.PostModel, err error)
	GetAllPost(ctx context.Context, username string, page int, limit int) (res []model.PostModel, err error)
	GetById(ctx context.Context, username string, id string) (res *model.PostModel, err error)
}
