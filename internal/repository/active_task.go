// internal/repository/active_task.go
package repository

import (
	"database/sql"
	"github.com/ze674/EZLine/internal/database"
)

// ActiveTaskRepository предоставляет методы для работы с активным заданием
type ActiveTaskRepository struct {
	db *sql.DB
}

// NewActiveTaskRepository создает новый репозиторий активного задания
func NewActiveTaskRepository() *ActiveTaskRepository {
	return &ActiveTaskRepository{
		db: database.DB,
	}
}

// SaveActiveTask сохраняет ID активного задания в базе данных
func (r *ActiveTaskRepository) SaveActiveTask(taskID int) error {
	// Сначала удаляем все существующие активные задания
	_, err := r.db.Exec("DELETE FROM active_task")
	if err != nil {
		return err
	}

	// Добавляем новое активное задание
	_, err = r.db.Exec(
		"INSERT INTO active_task (task_id) VALUES (?)",
		taskID)

	return err
}

// GetActiveTask возвращает ID активного задания из базы данных
func (r *ActiveTaskRepository) GetActiveTask() (int, error) {
	var taskID int

	err := r.db.QueryRow("SELECT task_id FROM active_task ORDER BY created_at DESC LIMIT 1").
		Scan(&taskID)

	if err == sql.ErrNoRows {
		return 0, nil // Нет активного задания
	}

	return taskID, err
}

// ClearActiveTask удаляет активное задание из базы данных
func (r *ActiveTaskRepository) ClearActiveTask() error {
	_, err := r.db.Exec("DELETE FROM active_task")
	return err
}
