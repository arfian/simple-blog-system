package server

import (
	"github.com/gin-gonic/gin"

	"simple-blog-system/internal/app/comment/port"
)

type (
	routes struct{}
)

var (
	Routes routes
)

func (r routes) New(router *gin.RouterGroup, handler port.ICommentHandler) {
	router.POST("/", handler.AddComment)
	router.PUT("/:id", handler.UpdateComment)
	router.DELETE("/:id", handler.DeleteComment)
	router.GET("/", handler.GetAllComment)
	router.GET("/:id", handler.GetCommentById)
}
