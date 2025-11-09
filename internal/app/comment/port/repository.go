package port

import (
	"context"
	"simple-blog-system/internal/app/comment/model"
)

type ICommentRepository interface {
	InsertComment(ctx context.Context, comment model.CommentModel) (model.CommentModel, error)
	UpdateComment(ctx context.Context, comment model.CommentModel) (res model.CommentModel, err error)
	DeleteComment(ctx context.Context, comment model.CommentModel) (err error)
	GetCommentById(ctx context.Context, id string) (res *model.CommentModel, err error)
	GetAllComment(ctx context.Context, page int, limit int) (res []model.CommentModel, err error)
}
