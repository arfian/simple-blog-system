package port

import (
	"github.com/gin-gonic/gin"
)

type ICommentHandler interface {

	// (POST /comment)
	AddComment(ctx *gin.Context)

	// (PUT /comment/:id)
	UpdateComment(ctx *gin.Context)

	// (DELETE /comment/:id)
	DeleteComment(ctx *gin.Context)

	// (GET /comment)
	GetAllComment(ctx *gin.Context)

	// (GET /comment/:id)
	GetCommentById(ctx *gin.Context)
}
