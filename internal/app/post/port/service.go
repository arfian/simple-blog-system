package port

import (
	"context"
	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/payload"
)

type IPostService interface {
	AddPost(ctx context.Context, username string, param payload.PostRequest) (res *model.PostModel, err error)
}
