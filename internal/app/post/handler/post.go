package handler

import (
	"simple-blog-system/internal/app/post/payload"
	"simple-blog-system/internal/app/post/port"
	"simple-blog-system/pkg/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	postService port.IPostService
}

func New(postService port.IPostService) port.IPostHandler {
	return &handler{
		postService: postService,
	}
}

func (h *handler) AddPost(c *gin.Context) {
	username := c.GetString("username")
	var (
		postRequest payload.PostRequest
	)

	if err := c.ShouldBind(&postRequest); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(postRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	res, err := h.postService.AddPost(c.Request.Context(), username, postRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "insert successfully",
		Data:    res,
	})
}

func (h *handler) UpdatePost(c *gin.Context) {
	username := c.GetString("username")
	var (
		postRequest payload.PostRequest
	)

	if err := c.ShouldBind(&postRequest); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(postRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	idStr := c.Param("id")

	res, err := h.postService.UpdatePost(c.Request.Context(), username, idStr, postRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "update successfully",
		Data:    res,
	})
}

func (h *handler) DeletePost(c *gin.Context) {
	username := c.GetString("username")
	var (
		postRequest payload.PostRequest
	)

	if err := c.ShouldBind(&postRequest); err != nil {
		helper.ResponseError(c, err)
		return
	}

	idStr := c.Param("id")

	res, err := h.postService.DeletePost(c.Request.Context(), username, idStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "delete successfully",
		Data:    res,
	})
}
