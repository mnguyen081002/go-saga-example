package dto

type CreateOrderRequest struct {
	ItemIDs []string `json:"item_ids"`
}

type UpdateItemRequest struct {
	Name string `json:"name"`
}
type SearchItemRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Paging
}
