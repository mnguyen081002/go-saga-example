package repository

import (
	"context"
	"fmt"
	"item-service/internal/dto"
	"item-service/internal/lib/db"
	"item-service/internal/models"
)

type SearchItemOption struct {
	dto.Paging
	Keyword string `json:"keyword"`
}

type UserRepository interface {
	Create(ctx context.Context, req models.User) (res models.User, err error)
	UpdateByGID(ctx context.Context, req models.User) (res models.User, err error)
	DB() *db.Database
}

type PaymentRepositoryImpl struct {
	db *db.Database
}

func (m *PaymentRepositoryImpl) DB() *db.Database {
	return m.db
}

func (m *PaymentRepositoryImpl) Create(ctx context.Context, req models.User) (res models.User, err error) {
	fmt.Println("req", req)
	err = m.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return models.User{}, err
	}
	return req, nil
}

func (m *PaymentRepositoryImpl) UpdateByGID(ctx context.Context, req models.User) (res models.User, err error) {
	err = m.db.WithContext(ctx).Model(&req).Where("gid = ?", req.GID).Updates(&req).Error
	if err != nil {
		return models.User{}, err
	}
	return req, nil
}

func NewItemRepository(engine *db.Database) UserRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &PaymentRepositoryImpl{
		db: engine,
	}

}
