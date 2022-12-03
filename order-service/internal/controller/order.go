package controller

import (
	"item-service/internal/dto"
	"item-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderController struct {
	OrderService service.OrderService
	logger      *zap.Logger
	BaseController
}

func InitOrderController(c *gin.RouterGroup, groupService service.OrderService, logger *zap.Logger) {
	controller := &OrderController{
		OrderService: groupService,
		logger:      logger,
	}
	g := c.Group("/order")
	g.POST("/", controller.Create)
}

func (b *OrderController) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	res, statusCode, err := b.OrderService.Create(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
