// internal/models/container.go
package models

import "time"

// internal/models/container.go
type Container struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	SerialNumber int       `json:"serial_number"` // Числовое поле
	TaskID       int       `json:"task_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
