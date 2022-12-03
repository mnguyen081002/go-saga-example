package repository

import (
	"context"
	"fmt"
	"item-service/internal/lib/db"
	"item-service/internal/models"

	"gorm.io/gorm"
)

type UserWalletRepository interface {
	Create(ctx context.Context, req models.UserWallet) (res models.UserWallet, err error)
	UpdateAmount(ctx context.Context, req models.UserWallet) (res models.UserWallet, err error)
	DB() *gorm.DB
}

type UserWalletRepositoryImpl struct {
	db *db.Database
}

func (m *UserWalletRepositoryImpl) DB() *gorm.DB {
	return m.db.DB
}

func (m *UserWalletRepositoryImpl) Create(ctx context.Context, req models.UserWallet) (res models.UserWallet, err error) {
	fmt.Println("req", req)
	err = m.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return models.UserWallet{}, err
	}
	return req, nil
}

func (m *UserWalletRepositoryImpl) UpdateAmount(ctx context.Context, req models.UserWallet) (res models.UserWallet, err error) {
	err = m.db.WithContext(ctx).Model(&req).Updates(&req).Error
	if err != nil {
		return models.UserWallet{}, err
	}
	return req, nil
}

func NewUserWalletRepository(engine *db.Database) UserWalletRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &UserWalletRepositoryImpl{
		db: engine,
	}
}
