package repository

import (
	"context"
	"item-service/internal/dto"
	"item-service/internal/lib/db"
	"item-service/internal/models"
)

type SearchItemOption struct {
	dto.Paging
	Keyword string `json:"keyword"`
}

type OrderRepository interface {
	Create(ctx context.Context, req models.Order) (res models.Order, err error)
}

type OrderRepositoryImpl struct {
	db *db.Database
}

func (m *OrderRepositoryImpl) Create(ctx context.Context, req models.Order) (res models.Order, err error) {
	err = m.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return models.Order{}, err
	}
	return req, nil
}

func NewItemRepository(engine *db.Database) OrderRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &OrderRepositoryImpl{db: engine}
}
