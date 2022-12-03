package service

import (
	"context"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"math/rand"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"gorm.io/gorm"
)

type (
	OrderService interface {
		Create(ctx context.Context, req dto.CreateOrderRequest) (item models.Order, statusCode int, err error)
	}

	OrderServiceImpl struct {
		db     *gorm.DB
		config config.Config
	}
)

func (s *OrderServiceImpl) Create(ctx context.Context, req dto.CreateOrderRequest) (item models.Order, statusCode int, err error) {
	// req := &gin.H{"amount": 30} // micro-service load threshold
	// DtmServer is the address of DTM micro-service
	// id order = random 7 string characters number
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	err = s.db.Create(models.Order{
		ID:      int8(r1.Intn(10000000)),
		ItemIDs: req.ItemIDs,
		Status:  "pending",
	}).Error

	if err != nil {
		return models.Order{}, 500, err
	}

	reqItem := map[string]interface{}{
		"quantity": 1,
		"item_id":  req.ItemIDs[0],
	}

	reqPayment := map[string]interface{}{
		"amount":  99999,
		"user_id": "1",
	}

	dtm := "http://localhost:36789/api/dtmsvr"
	saga := dtmcli.NewSaga(dtm, dtmcli.MustGenGid(dtm)).
		// Add("", "localhost:8080/api/order/compensation", req).
		Add(s.config.Services.ItemServiceURL+"/calculate-stock", s.config.Services.ItemServiceURL+"/compensation-stock", reqItem).
		Add(s.config.Services.PaymentServiceURL+"/purchase", s.config.Services.PaymentServiceURL+"/refund", reqPayment)
	saga.TimeoutToFail = 10
	saga.WithRetryLimit(1)
	err = saga.Submit()
	if err != nil {
		return item, 500, err
	}
	return
}

func NewItemService(db *gorm.DB, config config.Config) OrderService {
	return &OrderServiceImpl{
		db:     db,
		config: config,
	}
}
