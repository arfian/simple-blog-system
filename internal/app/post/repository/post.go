package repository

import (
	"context"

	"simple-blog-system/config/db"
	"simple-blog-system/pkg/cache"
	"simple-blog-system/pkg/transaction"

	"simple-blog-system/internal/app/post/model"
	"simple-blog-system/internal/app/post/port"
)

type repository struct {
	db    *db.GormDB
	cache cache.ICache
}

func NewRepository(db *db.GormDB) port.IPostRepository {
	return repository{db: db}
}

func (r repository) InsertPost(ctx context.Context, post model.PostModel) (model.PostModel, error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	qres := trx.Create(&post).Error

	return post, qres
}

func (r repository) UpdatePost(ctx context.Context, post model.PostModel) (res model.PostModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	qres := trx.Save(&post).Error

	return post, qres
}

func (r repository) GetPostById(ctx context.Context, id string) (res *model.PostModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Where("id = ?", id).First(&res).Error
	return res, err
}

func (r repository) DeletePost(ctx context.Context, post model.PostModel) (err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Delete(&post).Error
	return err
}

func (r repository) GetAllPost(ctx context.Context, page int, limit int) (res []model.PostModel, err error) {
	offset := (page - 1) * limit

	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}
