package common

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:now();not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:now();not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}
