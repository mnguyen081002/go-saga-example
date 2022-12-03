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
	logger       *zap.Logger
	BaseController
}

func InitOrderController(c *gin.RouterGroup, groupService service.OrderService, logger *zap.Logger) {
	controller := &OrderController{
		OrderService: groupService,
		logger:       logger,
	}
	g := c.Group("/order")
	g.POST("/create-order-saga", controller.CreateOrderSaga)
	g.POST("/create-order", controller.CreateOrder)
	g.POST("/cancel-order", controller.CancelOrder)
	g.POST("/order-success", controller.OrderSuccess)
}

func (b *OrderController) CreateOrderSaga(c *gin.Context) {
	var req dto.CreateOrderSagaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}

	res, statusCode, err := b.OrderService.CreateOrderSaga(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *OrderController) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderSagaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	req.GID, _ = c.GetQuery("gid")

	res, statusCode, err := b.OrderService.CreateOrder(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *OrderController) CancelOrder(c *gin.Context) {
	var req dto.CreateOrderSagaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	req.GID, _ = c.GetQuery("gid")

	res, statusCode, err := b.OrderService.CancelOrder(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *OrderController) OrderSuccess(c *gin.Context) {
	var req dto.CreateOrderSagaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	req.GID, _ = c.GetQuery("gid")

	res, statusCode, err := b.OrderService.OrderSuccess(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
