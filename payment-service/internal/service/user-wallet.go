package service

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"item-service/internal/repository"
)

type (
	UserWalletService interface {
		Create(ctx context.Context, req dto.CreateUserWalletRequest) (payment models.UserWallet, statusCode int, err error)
	}

	UserWalletServiceImpl struct {
		repo   repository.UserWalletRepository
		config config.Config
	}
)

func NewUserWalletService(repo repository.UserWalletRepository, config config.Config) UserWalletService {
	if repo == nil {
		panic("Item Repository is nil")
	}

	fmt.Println("repo", repo)

	return &UserWalletServiceImpl{
		repo:   repo,
		config: config,
	}
}

func (s *UserWalletServiceImpl) Create(ctx context.Context, req dto.CreateUserWalletRequest) (item models.UserWallet, statusCode int, err error) {

	item = models.UserWallet{
		UserID:  req.UserID,
		Balance: req.Balance,
	}

	if err := s.repo.DB().Create(&item).Error; err != nil {
		return models.UserWallet{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	return item, 200, err
}
