package port

import (
	"github.com/gin-gonic/gin"
)

type IPostHandler interface {

	// (POST /post)
	AddPost(ctx *gin.Context)

	// (PUT /post/:id)
	UpdatePost(ctx *gin.Context)

	// (DELETE /post/:id)
	DeletePost(ctx *gin.Context)

	// (GET /post)
	GetAllPost(ctx *gin.Context)

	// (GET /post/:id)
	GetById(ctx *gin.Context)
}
