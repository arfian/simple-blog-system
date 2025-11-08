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

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  helper.Response
// @Failure      400  {object}  helper.Response
// @Failure      404  {object}  helper.Response
// @Failure      500  {object}  helper.Response
// @Router       /accounts/{id} [get]
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
