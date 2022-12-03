package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dtm-labs/client/dtmcli"
	"gorm.io/gorm"
)

type (
	OrderService interface {
		CreateOrder(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error)
		CancelOrder(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error)
		CreateOrderSaga(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error)
		OrderSuccess(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error)
	}

	OrderServiceImpl struct {
		db     *gorm.DB
		config config.Config
	}
)

func NewItemService(db *gorm.DB, config config.Config) OrderService {
	return &OrderServiceImpl{
		db:     db,
		config: config,
	}
}

func (s *OrderServiceImpl) CreateOrderSaga(ctx context.Context, req dto.CreateOrderSagaRequest) (order models.Order, statusCode int, err error) {
	itemIDs := []string{}

	for _, orderItem := range req.OrderItems {
		itemIDs = append(itemIDs, orderItem.ItemID)
	}

	postBody, _ := json.Marshal(dto.GetItemsByIDsRequest{
		IDs: itemIDs,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(s.config.Services.ItemServiceURL+"/get-by-ids", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	res := dto.GetItemsByIDsResponse{}
	err = json.Unmarshal(body, &res)

	if err != nil {
		return models.Order{}, http.StatusInternalServerError, err
	}

	reqCreateTransaction := dto.CreateTransactionRequest{
		UserID: "1",
	}
	for _, item := range res.Data {
		reqCreateTransaction.Amount += item.Price
	}

	saga := dtmcli.NewSaga(s.config.Services.DTMServiceURL, dtmcli.MustGenGid(s.config.Services.DTMServiceURL)).
		Add(s.config.Server.MyURL+"/create-order", "cancel-order", req).
		Add(s.config.Services.ItemServiceURL+"/calculate-stock", s.config.Services.ItemServiceURL+"/compensation-stock", req).
		Add(s.config.Services.PaymentServiceURL+"/transaction", s.config.Services.PaymentServiceURL+"/refund", reqCreateTransaction).
		Add(s.config.Server.MyURL+"/order-success", "", req)
	saga.TimeoutToFail = 10000
	saga.WithRetryLimit(1)
	err = saga.Submit()
	if err != nil {
		return order, http.StatusInternalServerError, err
	}
	return
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	orderItems := []models.OrderItem{}
	for _, orderItem := range req.OrderItems {
		orderItems = append(orderItems, models.OrderItem{
			Quantity: orderItem.Quantity,
			ItemID:   orderItem.ItemID,
		})
	}

	if err := s.db.Create(&orderItems).Error; err != nil {
		return item, http.StatusInternalServerError, err
	}
	orderItemIDs := []string{}
	for _, orderItem := range orderItems {
		orderItemIDs = append(orderItemIDs, orderItem.ID.String())
	}

	err = s.db.Create(&models.Order{
		ID:           r1.Intn(100000000),
		OrderItemIDs: orderItemIDs,
		Status:       "pending",
		GID:          req.GID,
	}).Error
	if err != nil {
		return models.Order{}, http.StatusInternalServerError, err
	}

	return
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error) {
	err = s.db.Create(models.Order{
		Status: "canceled",
		GID:    req.GID,
	}).Error

	if err != nil {
		return models.Order{}, http.StatusInternalServerError, err
	}

	return
}

func (s *OrderServiceImpl) OrderSuccess(ctx context.Context, req dto.CreateOrderSagaRequest) (item models.Order, statusCode int, err error) {
	err = s.db.Create(&models.Order{
		Status: "success",
		GID:    req.GID,
	}).Error

	if err != nil {
		return models.Order{}, http.StatusInternalServerError, err
	}

	return
}
