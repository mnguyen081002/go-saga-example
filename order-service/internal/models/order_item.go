package models

import uuid "github.com/satori/go.uuid"

type OrderItem struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	OrderID  int8      `json:"order_id"`
	ItemID   string    `json:"item_id"`
	Quantity int       `json:"quantity"`
}
