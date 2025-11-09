package port

import (
	"context"
	"simple-blog-system/internal/app/post/model"
)

type IPostRepository interface {
	InsertPost(ctx context.Context, post model.PostModel) (model.PostModel, error)
	UpdatePost(ctx context.Context, post model.PostModel) (res model.PostModel, err error)
	DeletePost(ctx context.Context, post model.PostModel) (err error)
	GetPostById(ctx context.Context, id string) (res *model.PostModel, err error)
}
