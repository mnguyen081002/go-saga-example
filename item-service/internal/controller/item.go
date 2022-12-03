package controller

import (
	"item-service/internal/dto"
	"item-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ItemController struct {
	ItemService service.ItemService
	logger      *zap.Logger
	BaseController
}

func InitItemController(c *gin.RouterGroup, groupService service.ItemService, logger *zap.Logger) {
	controller := &ItemController{
		ItemService: groupService,
		logger:      logger,
	}
	g := c.Group("/item")
	g.GET("/:id", controller.GetByID)
	g.GET("/search", controller.Search)
	g.POST("/create", controller.Create)
	g.POST("/calculate-stock", controller.CalculateStock)
	g.POST("/compensation-stock", controller.CalculateStock)
	g.POST("/get-by-ids", controller.GetByIDs)
}

func (b *ItemController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	bg, err := b.ItemService.GetByID(c, idParam)
	if err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", bg)
}

func (b *ItemController) Search(c *gin.Context) {
	var req dto.SearchItemRequest

	if err := c.Bind(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}

	res, total, statusCode, err := b.ItemService.Search(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.ResponseList(c, "success", total, res)
}

func (b *ItemController) Create(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	res, statusCode, err := b.ItemService.Create(c, req)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}

func (b *ItemController) CalculateStock(c *gin.Context) {
	var req dto.CalculateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	statusCode, err := b.ItemService.CalculateStock(c, req.OrderItems)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", nil)
}

func (b *ItemController) CompensationStock(c *gin.Context) {
	var req dto.CalculateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	statusCode, err := b.ItemService.CompensationStock(c, req.OrderItems)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", nil)
}

func (b *ItemController) GetByIDs(c *gin.Context) {
	var req dto.GetByIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		b.ResponseError(c, http.StatusBadRequest, []error{err})
		return
	}
	res, statusCode, err := b.ItemService.GetByIDs(c, req.IDs)
	if err != nil {
		b.ResponseError(c, statusCode, []error{err})
		return
	}
	b.Response(c, http.StatusOK, "success", res)
}
