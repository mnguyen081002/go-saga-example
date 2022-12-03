package controller

import (
	"item-service/internal/dto"
	"item-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userService service.UserService
	logger      *zap.Logger
	BaseController
}

func InitUserController(c *gin.RouterGroup, userService service.UserService, logger *zap.Logger) {
	controller := &UserController{
		userService: userService,
		logger:      logger,
	}
	g := c.Group("/user")
	g.POST("/create", controller.Create)
	g.GET("/:id", controller.Get)
}

func (b *UserController) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	res, statusCode, err := b.userService.Create(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *UserController) Get(c *gin.Context) {
	id := c.Param("id")
	res, statusCode, err := b.userService.Get(c, id)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
