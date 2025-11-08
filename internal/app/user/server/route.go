package server

import (
	"github.com/gin-gonic/gin"

	"simple-blog-system/internal/app/user/port"
)

type (
	routes struct{}
)

var (
	Routes routes
)

func (r routes) New(router *gin.RouterGroup, handler port.IUserHandler) {
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
}

func (r routes) NewProfile(router *gin.RouterGroup, handler port.IUserHandler) {
	router.GET("/", handler.GetUser)
}
