package model

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type PostModel struct {
	ID        strfmt.UUID4 `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string       `json:"username" validate:"required"`
	Title     string       `json:"title" validate:"required"`
	Body      string       `json:"body" validate:"required"`
	Status    string       `json:"status"`
	CreatedBy string       `json:"created_by"`
	UpdatedBy string       `json:"updated_by" gorm:"default:null"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time   `json:"deleted_at" gorm:"default:null"`
}

func (u PostModel) TableName() string {
	return "posts"
}
