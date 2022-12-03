package service

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"

	"gorm.io/gorm"
)

type (
	UserWalletService interface {
		Create(ctx context.Context, req dto.CreateUserWalletRequest) (payment models.UserWallet, statusCode int, err error)
	}

	UserWalletServiceImpl struct {
		db     *gorm.DB
		config config.Config
	}
)

func NewUserWalletService(db *gorm.DB, config config.Config) UserWalletService {
	return &UserWalletServiceImpl{
		db:     db,
		config: config,
	}
}

func (s *UserWalletServiceImpl) Create(ctx context.Context, req dto.CreateUserWalletRequest) (item models.UserWallet, statusCode int, err error) {

	item = models.UserWallet{
		UserID:  req.UserID,
		Balance: req.Balance,
	}

	if err := s.db.Create(&item).Error; err != nil {
		return models.UserWallet{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	return item, 200, err
}
