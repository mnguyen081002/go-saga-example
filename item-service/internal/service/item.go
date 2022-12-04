package service

import (
	"context"
	"fmt"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"

	"net/http"

	"gorm.io/gorm"
)

type (
	ItemService interface {
		GetByID(ctx context.Context, id string) (item models.Item, err error)
		Search(ctx context.Context, req dto.SearchItemRequest) (res []models.Item, total *int64, statusCode int, err error)
		Create(ctx context.Context, req dto.CreateItemRequest) (item models.Item, statusCode int, err error)
		CalculateStock(ctx context.Context, irs []dto.OrderItem) (statusCode int, err error)
		CompensationStock(ctx context.Context, irs []dto.OrderItem) (statusCode int, err error)
		GetByIDs(ctx context.Context, ids []string) (res []models.Item, statusCode int, err error)
	}

	ItemServiceImpl struct {
		db     *gorm.DB
		config config.Config
	}
)

func NewItemService(db *gorm.DB, config config.Config) ItemService {
	return &ItemServiceImpl{
		db:     db,
		config: config,
	}
}

func (s *ItemServiceImpl) GetByID(ctx context.Context, id string) (item models.Item, err error) {
	err = s.db.First(&item, "id = ?", id).Error
	return
}

func (s *ItemServiceImpl) Search(ctx context.Context, req dto.SearchItemRequest) (res []models.Item, total *int64, statusCode int, err error) {
	fmt.Println(req)

	err = s.db.Find(res, "name LIKE ?", "%"+req.Keyword+"%").Count(total).Error

	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	return res, total, http.StatusOK, nil
}

func (s *ItemServiceImpl) CalculateStock(ctx context.Context, irs []dto.OrderItem) (statusCode int, err error) {
	for _, ir := range irs {

		item, err := s.GetByID(ctx, ir.ItemID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if int64(item.Stock) < ir.Quantity {
			return http.StatusBadRequest, fmt.Errorf("not enough stock")
		}

		stock := int64(item.Stock) - ir.Quantity
		if err := s.db.Model(models.Item{}).Where("id = ?", ir.ItemID).Update("stock", stock).Error; err != nil {
			return http.StatusInternalServerError, err
		}
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ItemServiceImpl) CompensationStock(ctx context.Context, irs []dto.OrderItem) (statusCode int, err error) {
	for _, ir := range irs {

		item, err := s.GetByID(ctx, ir.ItemID)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		stock := int64(item.Stock) + ir.Quantity
		if err := s.db.Model(models.Item{}).Update("stock", stock).Error; err != nil {
			return http.StatusInternalServerError, err
		}
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ItemServiceImpl) Create(ctx context.Context, req dto.CreateItemRequest) (item models.Item, statusCode int, err error) {
	err = s.db.Create(&models.Item{
		Name:         req.Name,
		Price:        req.Price,
		PriceMax:     req.PriceMax,
		PriceMin:     req.PriceMin,
		ShowFreeShip: req.ShowFreeShip,
		Description:  req.Description,
		SKU:          req.SKU,
		Quantity:     req.Quantity,
		Discount:     req.Discount,
		RawDiscount:  req.RawDiscount,
		Stock:        req.Stock,
		Images:       req.Images,
		CategoryID:   req.CategoryID,
		VariantIDs:   req.VariantIDs,
	}).Error
	return
}

func (s *ItemServiceImpl) GetByIDs(ctx context.Context, ids []string) (res []models.Item, statusCode int, err error) {
	err = s.db.Find(&res, "id IN ?", ids).Error
	return
}
