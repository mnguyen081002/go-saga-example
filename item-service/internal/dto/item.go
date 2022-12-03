package dto

type CreateItemRequest struct {
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	PriceMax float64 `json:"price_max"`
	PriceMin float64 `json:"price_min"`
	//PriceBeforeDiscount float64  `json:"price_before_discount"`
	ShowFreeShip bool     `json:"show_free_ship"`
	Description  string   `json:"description"`
	SKU          string   `json:"sku"`
	Quantity     int64    `json:"quantity"`
	Discount     string   `json:"discount"`
	RawDiscount  float64  `json:"raw_discount"`
	Stock        int64    `json:"stock"`
	Images       []string `json:"images"`
	CategoryID   string   `json:"category_id"`
	VariantIDs   []int64  `json:"variant_ids"`
}

type UpdateItemRequest struct {
	Name string `json:"name"`
}
type SearchItemRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Paging
}

type CalculateStockRequest struct {
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	Quantity int64  `json:"quantity"`
	ItemID   string `json:"item_id"`
}
