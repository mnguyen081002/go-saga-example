package service

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"item-service/internal/repository"

	"gorm.io/gorm"
)

type (
	PaymentService interface {
		CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (payment models.Transaction, statusCode int, err error)
		Refund(ctx context.Context, req dto.CreateTransactionRequest) (payment models.Transaction, statusCode int, err error)
	}

	PaymentServiceImpl struct {
		repo           repository.PaymentRepository
		userWalletRepo repository.UserWalletRepository
		config         config.Config
	}
)

func NewPaymentService(paymentRepo repository.PaymentRepository, userWalletRepo repository.UserWalletRepository, config config.Config) PaymentService {
	if paymentRepo == nil {
		panic("Item Repository is nil")
	}
	return &PaymentServiceImpl{
		repo:           paymentRepo,
		userWalletRepo: userWalletRepo,
		config:         config,
	}
}

func (s *PaymentServiceImpl) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (item models.Transaction, statusCode int, err error) {

	item = models.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "pending",
		GID:    req.GID,
	}

	if err := s.repo.DB().Create(&item).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	uw := models.UserWallet{}

	if err = s.userWalletRepo.DB().First(&uw, "user_id = ?", req.UserID).Error; err != nil {
		return models.Transaction{}, 500, err
	}

	if uw.Balance < req.Amount {
		return models.Transaction{}, 500, fmt.Errorf("Insufficient Balance")
	}

	if err := s.userWalletRepo.DB().Model(models.UserWallet{}).Where("user_id = ?", req.UserID).Update("balance", uw.Balance-req.Amount).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	item.Status = "success"
	if err := s.repo.DB().Updates(&item).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	return item, 200, err
}

func (s *PaymentServiceImpl) Refund(ctx context.Context, req dto.CreateTransactionRequest) (item models.Transaction, statusCode int, err error) {

	err = s.repo.DB().DB.First(&item, "g_id = ? and status = ?", req.GID, "pending").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Transaction{}, 200, fmt.Errorf("Transaction not found")
		}
		return models.Transaction{}, 200, err
	}

	userWallet := models.UserWallet{}
	err = s.userWalletRepo.DB().First(&userWallet, "user_id = ?", req.UserID).Error

	if err != nil {
		return models.Transaction{}, 500, err
	}
	userWallet.Balance = userWallet.Balance + req.Amount

	err = s.userWalletRepo.DB().Model(userWallet).Update("balance", userWallet.Balance).Where("user_id = ?", req.UserID).Error
	if err != nil {
		return models.Transaction{}, 500, err
	}

	item = models.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "refund",
		GID:    req.GID,
	}

	if err := s.repo.DB().Model(models.Transaction{}).Where("g_id = ? and status = ?", req.GID, "pending").Update("status", "refund").Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error update payment: %v", err)
	}

	return item, 200, err
}
