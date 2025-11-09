package handler

import (
	"simple-blog-system/internal/app/comment/payload"
	"simple-blog-system/internal/app/comment/port"
	"simple-blog-system/pkg/helper"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	commentService port.ICommentService
}

func New(commentService port.ICommentService) port.ICommentHandler {
	return &handler{
		commentService: commentService,
	}
}

func (h *handler) AddComment(c *gin.Context) {
	username := c.GetString("username")
	var (
		commentRequest payload.CommentRequest
	)

	if err := c.ShouldBind(&commentRequest); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(commentRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	res, err := h.commentService.AddComment(c.Request.Context(), username, commentRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "insert successfully",
		Data:    res,
	})
}

func (h *handler) UpdateComment(c *gin.Context) {
	username := c.GetString("username")
	var (
		commentRequest payload.CommentRequest
	)

	if err := c.ShouldBind(&commentRequest); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(commentRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	idStr := c.Param("id")

	res, err := h.commentService.UpdateComment(c.Request.Context(), username, idStr, commentRequest)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "update successfully",
		Data:    res,
	})
}

func (h *handler) DeleteComment(c *gin.Context) {
	username := c.GetString("username")

	idStr := c.Param("id")

	res, err := h.commentService.DeleteComment(c.Request.Context(), username, idStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "delete successfully",
		Data:    res,
	})
}

func (h *handler) GetAllComment(c *gin.Context) {
	username := c.GetString("username")

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	res, err := h.commentService.GetAllComment(c.Request.Context(), username, page, limit)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "get successfully",
		Data:    res,
	})
}

func (h *handler) GetCommentById(c *gin.Context) {
	username := c.GetString("username")

	idStr := c.Param("id")

	res, err := h.commentService.GetCommentById(c.Request.Context(), username, idStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "get successfully",
		Data:    res,
	})
}
