// internal/repository/container.go
package repository

import (
	"database/sql"
	"github.com/ze674/EZLine/internal/database"
	"github.com/ze674/EZLine/internal/models"
	"time"
)

const (
	StatusCreated = "created"
)

type ContainerRepository struct {
	db *sql.DB
}

func NewContainerRepository() *ContainerRepository {
	return &ContainerRepository{
		db: database.DB,
	}
}

// Создание контейнера с числовым серийным номером
func (r *ContainerRepository) CreateContainer(code string, serialNumber int, taskID int, status string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO containers (code, serial_number, task_id, status) VALUES (?, ?, ?, ?)",
		code, serialNumber, taskID, status)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Получение последнего серийного номера для задания
func (r *ContainerRepository) GetLastSerialNumber(taskID int) (int, error) {
	var serialNumber int

	err := r.db.QueryRow(
		"SELECT COALESCE(MAX(serial_number), 0) FROM containers WHERE task_id = ?",
		taskID).Scan(&serialNumber)

	if err != nil {
		return 0, err
	}

	return serialNumber, nil
}

// GetContainerByCode возвращает контейнер по коду
func (r *ContainerRepository) GetContainerByCode(code string) (*models.Container, error) {
	var container models.Container

	err := r.db.QueryRow(
		"SELECT id, code, task_id, status, created_at FROM containers WHERE code = ?",
		code).Scan(&container.ID, &container.Code, &container.TaskID, &container.Status, &container.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &container, nil
}

// GetContainersByTaskID возвращает контейнеры для задания
func (r *ContainerRepository) GetContainersByTaskID(taskID int) ([]models.Container, error) {
	rows, err := r.db.Query(
		"SELECT id, code, task_id, status, created_at FROM containers WHERE task_id = ? ORDER BY created_at DESC",
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []models.Container

	for rows.Next() {
		var container models.Container
		if err := rows.Scan(&container.ID, &container.Code, &container.TaskID, &container.Status, &container.CreatedAt); err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}

// UpdateContainerStatus обновляет статус контейнера
func (r *ContainerRepository) UpdateContainerStatus(id int64, status string) error {
	_, err := r.db.Exec(
		"UPDATE containers SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id)
	return err
}
