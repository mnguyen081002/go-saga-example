package dto

type UpdateItemRequest struct {
	Name string `json:"name"`
}
type SearchItemRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Paging
}

type CreateOrderSagaRequest struct {
	OrderItems []OrderItem `json:"order_items"`
	GID        string      `json:"gid"`
}

type OrderItem struct {
	Quantity int    `json:"quantity"`
	ItemID   string `json:"item_id"`
}

type CalculateStockRequest struct {
	OrderItems []OrderItem
}

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateTransactionRequest struct {
	Amount float64 `json:"amount"`
	UserID string  `json:"user_id"`
}

type GetItemsByIDsRequest struct {
	IDs []string `json:"ids"`
}

type GetItemsByIDsResponse struct {
	Data []Item `json:"data"`
}
