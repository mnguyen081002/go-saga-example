package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Transaction struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;default:now();not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp;default:now();not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
	UserID    string     `json:"user_id"`
	Amount    float64    `json:"amount"`
	Status    string     `json:"status"`
	GID       string     `json:"gid"`
}
