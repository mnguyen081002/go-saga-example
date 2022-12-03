package controller

import (
	"item-service/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (b *BaseController) Response(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, dto.SimpleResponse{
		Data:    data,
		Message: message,
	})
}

func (b *BaseController) ResponseList(c *gin.Context, message string, total *int64, data interface{}) {
	c.JSON(http.StatusOK, dto.SimpleResponseList{
		Message: message,
		Data:    data,
		Total:   total,
	})
}

func (b *BaseController) ResponseError(c *gin.Context, statusCode int, errs []error) {

	errorStrings := make([]string, len(errs))
	for i, err := range errs {
		c.Error(err)
		errorStrings[i] = err.Error()
	}

	c.AbortWithStatusJSON(statusCode, dto.ResponseError{
		Message: errs[0].Error(),
		Errors:  errorStrings,
	})
}
