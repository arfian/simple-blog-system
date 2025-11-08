package payload

import (
	"simple-blog-system/internal/app/user/model"
)

type User struct {
	User model.AuthUserModel `json:"auth_user"`
}
