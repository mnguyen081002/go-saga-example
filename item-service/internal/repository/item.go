package repository

import (
	"context"
	"fmt"
	"item-service/internal/dto"
	"item-service/internal/lib/db"
	"item-service/internal/models"
	"item-service/utils"

	"gorm.io/gorm"
)

type SearchItemOption struct {
	dto.Paging
	Keyword string `json:"keyword"`
}

type ItemRepository interface {
	GetByID(ctx context.Context, id string) (res models.Item, err error)
	Search(ctx context.Context, req SearchItemOption) (res []models.Item, total *int64, err error)
	Create(ctx context.Context, req models.Item) (res models.Item, err error)
	DB() *gorm.DB
}

type ItemRepositoryImpl struct {
	db *db.Database
}

func (m *ItemRepositoryImpl) DB() *gorm.DB {
	return m.db.DB
}

func (m *ItemRepositoryImpl) Create(ctx context.Context, req models.Item) (res models.Item, err error) {
	err = m.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return models.Item{}, err
	}
	return req, nil
}

func (m *ItemRepositoryImpl) Search(ctx context.Context, req SearchItemOption) (res []models.Item, total *int64, err error) {
	total = new(int64)
	err = m.db.Debug().
		Model(&models.Item{}).
		WithContext(ctx).
		Select("id,name,price,discount").
		Where("name ILIKE ?", "%"+req.Keyword+"%").
		Limit(req.Limit).Offset(req.Limit * (req.Page - 1)).
		Order(req.Sort).
		Count(total).
		Find(&res).Error
	if err != nil {
		return []models.Item{}, total, err
	}
	return
}

func NewItemRepository(engine *db.Database) ItemRepository {
	if engine == nil {
		panic("Database engine is null")
	}
	return &ItemRepositoryImpl{db: engine}
}

func (m *ItemRepositoryImpl) GetByID(ctx context.Context, id string) (res models.Item, err error) {
	err = m.db.WithContext(ctx).Where("id = ?", id).First(&res).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res, fmt.Errorf(utils.ItemNotFound)
		}
		return models.Item{}, err
	}
	return
}
