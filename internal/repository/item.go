// internal/repository/item.go
package repository

import (
	"database/sql"
	"github.com/ze674/EZLine/internal/database"
	"github.com/ze674/EZLine/internal/models"
)

const (
	StatusAggregated = "aggregated"
	StatusScanned    = "scanned"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		db: database.DB,
	}
}

// CreateItem создает новую запись о товаре
func (r *ItemRepository) CreateItem(code string, taskID int, status string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO items (code, task_id, status) VALUES (?, ?, ?)",
		code, taskID, status)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *ItemRepository) CreateItems(codes []string, taskID int, status string) error {
	for _, code := range codes {
		_, err := r.CreateItem(code, taskID, status)
		if err != nil {
			return err
		}
	}
	return nil
}

// AssignItemToContainer привязывает товар к контейнеру
func (r *ItemRepository) AssignItemToContainer(itemID, containerID int64) error {
	_, err := r.db.Exec(
		"UPDATE items SET container_id = ?, status = ? WHERE id = ?",
		containerID, StatusAggregated, itemID)
	return err
}

// GetItemsByTaskID возвращает товары для задания
func (r *ItemRepository) GetItemsByTaskID(taskID int) ([]models.Item, error) {
	rows, err := r.db.Query(
		"SELECT id, code, task_id, container_id, status, created_at FROM items WHERE task_id = ? ORDER BY created_at DESC",
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Code, &item.TaskID, &item.ContainerID, &item.Status, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetItemsByContainerID возвращает товары для контейнера
func (r *ItemRepository) GetItemsByContainerID(containerID int64) ([]models.Item, error) {
	rows, err := r.db.Query(
		"SELECT id, code, task_id, container_id, status, created_at FROM items WHERE container_id = ? ORDER BY created_at DESC",
		containerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Code, &item.TaskID, &item.ContainerID, &item.Status, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetItemByCode возвращает товар по коду
func (r *ItemRepository) GetItemByCode(code string) (*models.Item, error) {
	var item models.Item

	err := r.db.QueryRow(
		"SELECT id, code, task_id, container_id, status, created_at FROM items WHERE code = ?",
		code).Scan(&item.ID, &item.Code, &item.TaskID, &item.ContainerID, &item.Status, &item.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// UpdateItemStatus обновляет статус товара
func (r *ItemRepository) UpdateItemStatus(id int64, status string) error {
	_, err := r.db.Exec(
		"UPDATE items SET status = ? WHERE id = ?",
		status, id)
	return err
}
