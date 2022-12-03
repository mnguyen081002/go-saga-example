package dto

type CreateTransactionRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	GID   string  `json:"gid"`
}

type CreateUserWalletRequest struct {
	UserID string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type UpdateItemRequest struct {
	Name string `json:"name"`
}
type SearchItemRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Paging
}
