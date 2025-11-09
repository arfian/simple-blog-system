package handler

import (
	"simple-blog-system/internal/app/post/payload"
	"simple-blog-system/internal/app/post/port"
	"simple-blog-system/pkg/helper"
	"strconv"

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

// @BasePath /v1

// @Summary Add Post
// @Description Add Post
// @Tags post
// @Accept json
// @Produce json
// @Param post body payload.PostRequest true "Param Post"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/post [post]
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

// @Summary Update Post
// @Description Update Post
// @Tags post
// @Accept json
// @Produce json
// @Param post body payload.PostRequest true "Param Post"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/post/{id} [put]
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

// @Summary Delete Post
// @Description Delete Post
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/post/{id} [delete]
func (h *handler) DeletePost(c *gin.Context) {
	username := c.GetString("username")

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

// @Summary Get All Post
// @Description Get All Post
// @Tags post
// @Accept json
// @Produce json
// @Param page path int true "Page"
// @Param limit path int true "Limit"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/post [get]
func (h *handler) GetAllPost(c *gin.Context) {
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

	res, err := h.postService.GetAllPost(c.Request.Context(), username, page, limit)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "get successfully",
		Data:    res,
	})
}

// @Summary Get Post ID
// @Description Get Post ID
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/post/{id} [get]
func (h *handler) GetById(c *gin.Context) {
	username := c.GetString("username")

	idStr := c.Param("id")

	res, err := h.postService.GetById(c.Request.Context(), username, idStr)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "get successfully",
		Data:    res,
	})
}
