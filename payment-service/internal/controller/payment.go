package controller

import (
	"item-service/internal/dto"
	"item-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentController struct {
	paymentService service.PaymentService
	logger         *zap.Logger
	BaseController
}

func InitPaymentController(c *gin.RouterGroup, groupService service.PaymentService, logger *zap.Logger) {
	controller := &PaymentController{
		paymentService: groupService,
		logger:         logger,
	}
	g := c.Group("/payment")
	g.POST("/transaction", controller.Create)
	g.POST("/refund", controller.Refund)
}

func (b *PaymentController) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	req.GID, _ = c.GetQuery("gid")

	res, statusCode, err := b.paymentService.CreateTransaction(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *PaymentController) Refund(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	req.GID, _ = c.GetQuery("gid")

	res, statusCode, err := b.paymentService.Refund(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
