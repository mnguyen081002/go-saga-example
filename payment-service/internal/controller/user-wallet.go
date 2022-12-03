package controller

import (
	"item-service/internal/dto"
	"item-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserWalletController struct {
	UserWalletService service.UserWalletService
	logger            *zap.Logger
	BaseController
}

func InitUserWalletController(c *gin.RouterGroup, groupService service.UserWalletService, logger *zap.Logger) {
	controller := &UserWalletController{
		UserWalletService: groupService,
		logger:            logger,
	}
	g := c.Group("/user-wallet")
	g.POST("/", controller.Create)
}

func (b *UserWalletController) Create(c *gin.Context) {
	var req dto.CreateUserWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	res, statusCode, err := b.UserWalletService.Create(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
