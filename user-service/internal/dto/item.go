package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
}

type UpdateItemRequest struct {
	Name string `json:"name"`
}
type SearchItemRequest struct {
	Keyword string `json:"keyword" form:"keyword"`
	Paging
}
