package models

import (
	"time"

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
	ID        int8           `json:"id" gorm:"primary_key"`
	ItemIDs   pq.StringArray `json:"item_ids" gorm:"type:text[]"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;default:now();not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;default:now();not null" json:"updated_at"`
	DeletedAt *time.Time     `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}
