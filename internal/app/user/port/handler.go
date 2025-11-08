package port

import (
	"github.com/gin-gonic/gin"
)

type IUserHandler interface {

	// (POST /user/register)
	Register(ctx *gin.Context)

	// (POST /user/login)
	Login(ctx *gin.Context)

	// (GET /user/)
	GetUser(ctx *gin.Context)
}
