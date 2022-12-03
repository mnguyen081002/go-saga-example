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

type PaymentRepository interface {
	Create(ctx context.Context, req models.Transaction) (res models.Transaction, err error)
	UpdateByGID(ctx context.Context, req models.Transaction) (res models.Transaction, err error)
	DB() *db.Database
}

type PaymentRepositoryImpl struct {
	db *db.Database
}

func (m *PaymentRepositoryImpl) DB() *db.Database {
	return m.db
}

func (m *PaymentRepositoryImpl) Create(ctx context.Context, req models.Transaction) (res models.Transaction, err error) {

	// random string generator
	err = m.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return models.Transaction{}, err
	}
	return req, nil
}

func (m *PaymentRepositoryImpl) UpdateByGID(ctx context.Context, req models.Transaction) (res models.Transaction, err error) {
	err = m.db.WithContext(ctx).Model(&req).Where("gid = ?", req.GID).Updates(&req).Error
	if err != nil {
		return models.Transaction{}, err
	}
	return req, nil
}

func NewPaymentRepository(engine *db.Database) PaymentRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &PaymentRepositoryImpl{
		db: engine,
	}

}
