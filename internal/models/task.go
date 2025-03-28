package models

import "time"

// Статусы заданий
const (
	TaskStatusNew        = "новое"
	TaskStatusInProgress = "в работе"
	TaskStatusCompleted  = "завершено"
)

// Task представляет производственное задание
type Task struct {
	ID          int       `json:"ID"`          // Уникальный идентификатор
	ProductID   int       `json:"ProductID"`   // ID продукта
	ProductName string    `json:"ProductName"` // Название продукта
	LineID      int       `json:"LineID"`      // ID производственной линии
	LineName    string    `json:"LineName"`    // Название линии
	Date        string    `json:"Date"`        // Дата в формате ДД.ММ.ГГГГ
	BatchNumber string    `json:"BatchNumber"` // Номер партии
	Status      string    `json:"Status"`      // Статус задания
	CreatedAt   time.Time `json:"CreatedAt"`   // Дата создания
}
