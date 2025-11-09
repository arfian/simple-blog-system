package server

import (
	"github.com/gin-gonic/gin"

	"simple-blog-system/internal/app/post/port"
)

type (
	routes struct{}
)

var (
	Routes routes
)

func (r routes) New(router *gin.RouterGroup, handler port.IPostHandler) {
	router.POST("/", handler.AddPost)
	router.PUT("/:id", handler.UpdatePost)
	router.DELETE("/:id", handler.DeletePost)
	router.GET("/", handler.GetAllPost)
	router.GET("/:id", handler.GetById)
}
