package service

import (
	"context"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"item-service/internal/repository"
)

type (
	UserService interface {
		Create(ctx context.Context, req dto.CreateUserRequest) (payment models.User, statusCode int, err error)
		Get(ctx context.Context, id string) (payment models.User, statusCode int, err error)
	}

	PaymentServiceImpl struct {
		repo   repository.UserRepository
		config config.Config
	}
)

func (s *PaymentServiceImpl) Create(ctx context.Context, req dto.CreateUserRequest) (item models.User, statusCode int, err error) {
	err = s.repo.DB().DB.First(&item, "gid = ?", req.GID).Error

	if err != nil {
		return models.User{}, 500, err
	}

	return item, 200, err
}

func (s *PaymentServiceImpl) Get(ctx context.Context, id string) (item models.User, statusCode int, err error) {
	err = s.repo.DB().DB.First(&item, "id = ?", id).Error
	if err != nil {
		return models.User{}, 500, err
	}

	return item, 200, err
}

func NewItemService(itemRepo repository.UserRepository, config config.Config) UserService {
	if itemRepo == nil {
		panic("Item Repository is nil")
	}
	return &PaymentServiceImpl{
		repo:   itemRepo,
		config: config,
	}
}
