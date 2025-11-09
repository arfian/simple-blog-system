package model

import (
	"time"

	model "simple-blog-system/internal/app/post/model"

	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type CommentModel struct {
	ID        strfmt.UUID4 `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Username  string       `json:"username" validate:"required"`
	Comment   string       `json:"comment" validate:"required"`
	PostId    string       `json:"post_id" validate:"required"`
	Post      model.PostModel
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by" gorm:"default:null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"default:null"`
}

func (u CommentModel) TableName() string {
	return "comments"
}
