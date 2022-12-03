package service

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"item-service/internal/repository"
	"math/rand"
	"time"

	"github.com/dtm-labs/client/dtmcli"
)

type (
	OrderService interface {
		Create(ctx context.Context, req dto.CreateOrderRequest) (item models.Order, statusCode int, err error)
	}

	OrderServiceImpl struct {
		repo   repository.OrderRepository
		config config.Config
	}
)

func (s *OrderServiceImpl) Create(ctx context.Context, req dto.CreateOrderRequest) (item models.Order, statusCode int, err error) {
	// req := &gin.H{"amount": 30} // micro-service load threshold
	// DtmServer is the address of DTM micro-service
	// id order = random 7 string characters number
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	fmt.Print()
	s.repo.Create(ctx, models.Order{
		ID:      int8(r1.Intn(10000000)),
		ItemIDs: req.ItemIDs,
		Status:  "pending",
	})

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

func NewItemService(itemRepo repository.OrderRepository, config config.Config) OrderService {
	if itemRepo == nil {
		panic("Item Repository is nil")
	}
	return &OrderServiceImpl{
		repo:   itemRepo,
		config: config,
	}
}
