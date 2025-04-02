// internal/models/item.go
package models

import "time"

// Item представляет единицу продукции с кодом маркировки
type Item struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	TaskID      int       `json:"task_id"`
	ContainerID *int64    `json:"container_id"` // Может быть NULL
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
