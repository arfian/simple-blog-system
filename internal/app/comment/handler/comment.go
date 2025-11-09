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

// @BasePath /v1

// @Summary Add Comment
// @Description Add Comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment body payload.CommentRequest true "Param Comment"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/comment [post]
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

// @Summary Update Comment
// @Description Update Comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment body payload.CommentRequest true "Param Comment"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/comment/{id} [put]
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

// @Summary Delete Comment
// @Description Delete Comment
// @Tags comment
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/comment/{id} [delete]
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

// @Summary Get All Comment
// @Description Get All Comment
// @Tags comment
// @Accept json
// @Produce json
// @Param page path int true "Page"
// @Param limit path int true "Limit"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/comment [get]
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

// @Summary Get Comment ID
// @Description Get Comment ID
// @Tags comment
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/comment/{id} [get]
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
