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
	PaymentService interface {
		CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (payment models.Transaction, statusCode int, err error)
		Refund(ctx context.Context, req dto.CreateTransactionRequest) (payment models.Transaction, statusCode int, err error)
	}

	PaymentServiceImpl struct {
		db     *gorm.DB
		config config.Config
	}
)

func NewPaymentService(db *gorm.DB, config config.Config) PaymentService {
	return &PaymentServiceImpl{
		db:     db,
		config: config,
	}
}

func (s *PaymentServiceImpl) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (item models.Transaction, statusCode int, err error) {

	item = models.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "pending",
		GID:    req.GID,
	}

	fmt.Println(req.GID)

	if err := s.db.Create(&item).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	uw := models.UserWallet{}

	if err = s.db.First(&uw, "user_id = ?", req.UserID).Error; err != nil {
		return models.Transaction{}, 500, err
	}

	if uw.Balance < req.Amount {
		return models.Transaction{}, 500, fmt.Errorf("Insufficient Balance")
	}

	if err := s.db.Model(models.UserWallet{}).Where("user_id = ?", req.UserID).Update("balance", uw.Balance-req.Amount).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	item.Status = "success"
	if err := s.db.Model(models.Transaction{}).Where("g_id = ?", req.GID).Updates(&item).Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error creating payment: %v", err)
	}

	return item, 200, nil
}

func (s *PaymentServiceImpl) Refund(ctx context.Context, req dto.CreateTransactionRequest) (item models.Transaction, statusCode int, err error) {

	err = s.db.First(&item, "g_id = ? and status = ?", req.GID, "success").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Transaction{}, 200, fmt.Errorf("Transaction not found")
		}
		return models.Transaction{}, 200, err
	}

	userWallet := models.UserWallet{}
	err = s.db.First(&userWallet, "user_id = ?", req.UserID).Error

	if err != nil {
		return models.Transaction{}, 500, err
	}
	userWallet.Balance = userWallet.Balance + req.Amount

	err = s.db.Model(userWallet).Update("balance", userWallet.Balance).Where("user_id = ?", req.UserID).Error
	if err != nil {
		return models.Transaction{}, 500, err
	}

	item = models.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "refund",
		GID:    req.GID,
	}

	if err := s.db.Model(models.Transaction{}).Where("g_id = ? and status = ?", req.GID, "pending").Update("status", "refund").Error; err != nil {
		return models.Transaction{}, 500, fmt.Errorf("error update payment: %v", err)
	}

	return item, 200, err
}
