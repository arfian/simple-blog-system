package repository

import (
	"context"

	"simple-blog-system/config/db"
	"simple-blog-system/pkg/cache"
	"simple-blog-system/pkg/transaction"

	"simple-blog-system/internal/app/comment/model"
	"simple-blog-system/internal/app/comment/port"
)

type repository struct {
	db    *db.GormDB
	cache cache.ICache
}

func NewRepository(db *db.GormDB) port.ICommentRepository {
	return repository{db: db}
}

func (r repository) InsertComment(ctx context.Context, post model.CommentModel) (model.CommentModel, error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	qres := trx.Create(&post).Error

	return post, qres
}

func (r repository) UpdateComment(ctx context.Context, post model.CommentModel) (res model.CommentModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	qres := trx.Save(&post).Error

	return post, qres
}

func (r repository) GetCommentById(ctx context.Context, id string) (res *model.CommentModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Preload("Post").Where("id = ?", id).First(&res).Error
	return res, err
}

func (r repository) DeleteComment(ctx context.Context, comment model.CommentModel) (err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Delete(&comment).Error
	return err
}

func (r repository) GetAllComment(ctx context.Context, page int, limit int) (res []model.CommentModel, err error) {
	offset := (page - 1) * limit

	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Preload("Post").Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}
