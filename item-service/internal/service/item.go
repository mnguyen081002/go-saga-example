package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"item-service/config"
	"item-service/internal/dto"
	"item-service/internal/models"
	"item-service/internal/repository"
	"log"
	"net/http"
)

type (
	ItemService interface {
		GetByID(ctx context.Context, id string) (item models.Item, err error)
		Search(ctx context.Context, req dto.SearchItemRequest) (res []models.Item, total *int64, statusCode int, err error)
		Create(ctx context.Context, req dto.CreateItemRequest) (item models.Item, statusCode int, err error)
		CalculateStock(ctx context.Context, itemID string, quantity int64) (statusCode int, err error)
		CompensationStock(ctx context.Context, itemID string, quantity int64) (statusCode int, err error)
	}

	ItemServiceImpl struct {
		repo   repository.ItemRepository
		config config.Config
	}
)

func getInfoFromOtherServiceById(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	r := dto.ResponseError{}
	json.Unmarshal(body, &r)

	fmt.Println(r)

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf(r.Errors[0])
	}
	return resp.StatusCode, nil
}

func getInfoFromOtherServiceByIds(url string, ids []int64) (int, error) {
	body := dto.GetByIDsRequest{
		IDs: ids,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r := struct {
		Errors  []string `json:"errors"`
		Total   int64    `json:"total"`
		Message string   `json:"message"`
	}{}
	json.Unmarshal(bodyResp, &r)

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf(r.Errors[0])
	}

	if r.Total != int64(len(ids)) {
		return http.StatusBadRequest, fmt.Errorf("invalid variant ids")
	}

	return resp.StatusCode, nil
}

func (s *ItemServiceImpl) Create(ctx context.Context, req dto.CreateItemRequest) (item models.Item, statusCode int, err error) {
	if code, err := getInfoFromOtherServiceById(s.config.Services.CategoryServiceURL + req.CategoryID); err != nil {
		return models.Item{}, code, err
	}

	if code, err := getInfoFromOtherServiceByIds(s.config.Services.VariantServiceURL+"ids", req.VariantIDs); err != nil {
		return models.Item{}, code, err
	}

	item, err = s.repo.Create(ctx, models.Item{
		Name:     req.Name,
		Price:    req.Price,
		PriceMax: req.PriceMax,
		PriceMin: req.PriceMin,
		//PriceBeforeDiscount: req.PriceBeforeDiscount,
		ShowFreeShip: req.ShowFreeShip,
		Description:  req.Description,
		SKU:          req.SKU,
		Quantity:     req.Quantity,
		Discount:     req.Discount,
		RawDiscount:  req.RawDiscount,
		//Stock:               req.Stock,
		Images:     req.Images,
		CategoryID: req.CategoryID,
		VariantIDs: req.VariantIDs,
	})
	return
}

func (s *ItemServiceImpl) Search(ctx context.Context, req dto.SearchItemRequest) (res []models.Item, total *int64, statusCode int, err error) {
	fmt.Println(req)

	res, total, err = s.repo.Search(ctx, repository.SearchItemOption{
		Keyword: req.Keyword,
		Paging: dto.Paging{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		},
	})

	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	return res, total, http.StatusOK, nil
}

func NewItemService(itemRepo repository.ItemRepository, config config.Config) ItemService {
	if itemRepo == nil {
		panic("Item Repository is nil")
	}
	return &ItemServiceImpl{
		repo:   itemRepo,
		config: config,
	}
}

func (s *ItemServiceImpl) GetByID(ctx context.Context, id string) (item models.Item, err error) {
	item, err = s.repo.GetByID(ctx, id)
	return
}

func (s *ItemServiceImpl) CalculateStock(ctx context.Context, itemID string, quantity int64) (statusCode int, err error) {
	item, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if int64(item.Stock) < quantity {
		return http.StatusBadRequest, fmt.Errorf("not enough stock")
	}

	stock := int64(item.Stock) - quantity

	if err := s.repo.DB().Model(models.Item{}).Update("stock", stock).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *ItemServiceImpl) CompensationStock(ctx context.Context, itemID string, quantity int64) (statusCode int, err error) {
	item, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	stock := int64(item.Stock) + quantity

	if err := s.repo.DB().Model(models.Item{}).Update("stock", stock).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
