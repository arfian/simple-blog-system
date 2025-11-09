package port

import (
	"github.com/gin-gonic/gin"
)

type IPostHandler interface {

	// (POST /post)
	AddPost(ctx *gin.Context)

	// (POST /put)
	UpdatePost(ctx *gin.Context)

	// (DELETE /put)
	DeletePost(ctx *gin.Context)
}
