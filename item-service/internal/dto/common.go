package dto

type ResponseError struct {
	//Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

type SimpleResponse struct {
	//Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SimpleResponseList struct {
	//Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   *int64      `json:"total"`
}

type Paging struct {
	Page  int    `form:"page" json:"page" binding:"required"`
	Limit int    `form:"limit" json:"limit" binding:"required"`
	Sort  string `form:"sort" json:"sort" binding:"required"`
}

type GetByIDsRequest struct {
	IDs []string `json:"ids"`
}
