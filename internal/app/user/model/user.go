package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type AuthUserModel struct {
	ID        strfmt.UUID4 `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string       `json:"username" validate:"required"`
	Password  string       `json:"password" validate:"required"`
	IsActive  bool         `json:"is_active"`
	LastLogin time.Time    `json:"last_login"`
	CreatedBy string       `json:"created_by"`
	UpdatedBy string       `json:"updated_by" gorm:"default:null"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time    `json:"deleted_at" gorm:"default:null"`
}

func (u AuthUserModel) TableName() string {
	return "auth_user"
}
