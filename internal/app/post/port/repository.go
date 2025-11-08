package port

import (
	"context"
	"simple-blog-system/internal/app/post/model"
)

type IPostRepository interface {
	InsertPost(ctx context.Context, post model.PostModel) (model.PostModel, error)
	GetPostById(ctx context.Context, year int, month int, id string) (res *model.PostModel, err error)
}
