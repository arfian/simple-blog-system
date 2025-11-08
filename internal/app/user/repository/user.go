package repository

import (
	"context"

	"simple-blog-system/config/db"
	"simple-blog-system/pkg/cache"
	"simple-blog-system/pkg/transaction"

	"simple-blog-system/internal/app/user/model"
	"simple-blog-system/internal/app/user/port"
)

type repository struct {
	db    *db.GormDB
	cache cache.ICache
}

func NewRepository(db *db.GormDB) port.IUserRepository {
	return repository{db: db}
}

func (r repository) InsertUser(ctx context.Context, user model.AuthUserModel) (model.AuthUserModel, error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	qres := trx.Create(&user).Error

	return user, qres
}

func (r repository) GetUserByUsername(ctx context.Context, username string) (user []model.AuthUserModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Select("id, username, created_at, updated_at").Where("username = ?", username).Find(&user).Error
	return user, err
}

func (r repository) GetPasswordByUsername(ctx context.Context, username string) (user []model.AuthUserModel, err error) {
	trx := transaction.GetTrxContext(ctx, r.db)
	err = trx.Select("id, password, username, created_at, updated_at").Where("username = ?", username).Find(&user).Error
	return user, err
}

func (r repository) UpdateLastLogin(ctx context.Context, user model.AuthUserModel) error {
	trx := transaction.GetTrxContext(ctx, r.db)
	err := trx.Model(&model.AuthUserModel{}).Where("username = ?", user.Username).Update("last_login", user.LastLogin).Error
	return err
}
