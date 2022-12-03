package models

import (
	"item-service/internal/models/common"

	"github.com/lib/pq"
)

type Rating struct {
	RatingStar  float64 `json:"rating_star"`
	RatingCount []int   `json:"rating_count"`
}

type TierVariation struct {
	Images      []string       `json:"images"`
	Name        string         `json:"name"`
	Options     pq.StringArray `json:"options"`
	SummedStock int            `json:"summed_stock"`
}

type Item struct {
	common.BaseModel
	Name                string  `json:"name"`
	TSV                 string  `gorm:"type:tsvector" json:"-"`
	Price               float64 `json:"price"`
	PriceMax            float64 `json:"price_max"`
	PriceMin            float64 `json:"price_min"`
	PriceBeforeDiscount float64 `json:"price_before_discount"`
	ShowFreeShip        bool    `json:"show_free_ship"`
	// TODO: Think about how to implement this
	//Sold int `json:"sold"`

	Description string         `json:"description"`
	SKU         string         `json:"sku" gorm:"unique"`
	Quantity    int64          `json:"quantity"`
	Discount    string         `json:"discount"`
	RawDiscount float64        `json:"raw_discount"`
	Stock       int64          `json:"stock"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	VariantIDs  pq.Int64Array  `gorm:"type:integer[]"`

	Liked      bool `json:"liked" gorm:"-"`
	LikedCount int  `json:"liked_count"`
	//Brand       string   `json:"brand"`
	//BrandID     int      `json:"brand_id"`
	CategoryID     string            `json:"category_id"`
	Categories     []common.Category `json:"categories" gorm:"-"`
	TierVariations []TierVariation   `json:"tier_variations" gorm:"-"`
	//
	Rating Rating `json:"rating" gorm:"-"`
}
