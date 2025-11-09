package port

import (
	"context"
	"simple-blog-system/internal/app/comment/model"
	"simple-blog-system/internal/app/comment/payload"
)

type ICommentService interface {
	AddComment(ctx context.Context, username string, param payload.CommentRequest) (res *model.CommentModel, err error)
	UpdateComment(ctx context.Context, username string, id string, param payload.CommentRequest) (res *model.CommentModel, err error)
	DeleteComment(ctx context.Context, username string, id string) (res *model.CommentModel, err error)
	GetAllComment(ctx context.Context, username string, page int, limit int) (res []model.CommentModel, err error)
	GetCommentById(ctx context.Context, username string, id string) (res *model.CommentModel, err error)
}
