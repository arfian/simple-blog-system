package port

import (
	"github.com/gin-gonic/gin"
)

type IPostHandler interface {

	// (POST /post)
	AddPost(ctx *gin.Context)
}
