package handler

import (
	"simple-blog-system/internal/app/user/model"
	"simple-blog-system/internal/app/user/payload"
	"simple-blog-system/internal/app/user/port"
	"simple-blog-system/pkg/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	userService port.IUserService
}

func New(userService port.IUserService) port.IUserHandler {
	return &handler{
		userService: userService,
	}
}

// @BasePath /v1

// @Summary Register User
// @Description Register User
// @Tags user
// @Accept json
// @Produce json
// @Param user body payload.User true "Param Register"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /public-api/user/register [post]
func (h *handler) Register(c *gin.Context) {
	var (
		dataUser payload.User
	)
	if err := c.ShouldBind(&dataUser); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(dataUser)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	res, err := h.userService.Register(c.Request.Context(), dataUser.User)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "register successfully",
		Data:    res,
	})
}

// @Summary Login User
// @Description Login User
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.AuthUserModel true "Param Login"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /public-api/user/login [post]
func (h *handler) Login(c *gin.Context) {
	var (
		dataUser model.AuthUserModel
	)

	if err := c.ShouldBind(&dataUser); err != nil {
		helper.ResponseError(c, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(dataUser)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	res, err := h.userService.Login(c.Request.Context(), dataUser)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "login successfully",
		Data:    res,
	})
}

// @Summary Get User
// @Description Get User
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /api/profile [get]
func (h *handler) GetUser(c *gin.Context) {
	username := c.GetString("username")
	res, err := h.userService.GetUser(c.Request.Context(), username)
	if err != nil {
		helper.ResponseError(c, err)
		return
	}

	helper.ResponseData(c, &helper.Response{
		Message: "get user successfully",
		Data:    res,
	})
}
