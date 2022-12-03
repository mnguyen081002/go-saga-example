package models

import (
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

type Order struct {
	ID           int            `json:"id" gorm:"primary_key"`
	GID          string         `json:"gid"`
	OrderItemIDs pq.StringArray `json:"item_ids" gorm:"type:text[]"`
	Status       string         `json:"status"`
}
