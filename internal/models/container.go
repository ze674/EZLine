// internal/models/container.go
package models

import "time"

// Container представляет групповую упаковку (короб, паллету и т.д.)
type Container struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	TaskID    int       `json:"task_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
