// internal/services/task_service.go
package services

import (
	"fmt"
	"github.com/ze674/EZLine/internal/api"
	"github.com/ze674/EZLine/internal/models"
	"sync"
)

// TaskService предоставляет методы для работы с заданиями
type TaskService struct {
	mu            sync.Mutex
	factoryClient *api.FactoryClient
	lineID        int
	activeTaskID  int // Тут храним ID активного задания
}

// NewTaskService создает новый сервис для управления заданиями
func NewTaskService(factoryClient *api.FactoryClient, lineID int) *TaskService {
	return &TaskService{
		factoryClient: factoryClient,
		lineID:        lineID,
		activeTaskID:  0, // Изначально нет активного задания
	}
}

// GetLineID возвращает ID производственной линии
func (s *TaskService) GetLineID() int {
	return s.lineID
}

// GetTasks получает список заданий с сервера
func (s *TaskService) GetTasks() ([]models.Task, error) {
	return s.factoryClient.GetTasks(s.lineID)
}

// GetTaskByID получает информацию о задании по ID
func (s *TaskService) GetTaskByID(taskID int) (models.Task, error) {
	return s.factoryClient.GetTaskByID(taskID)
}

// GetActiveTaskID возвращает ID активного задания
func (s *TaskService) GetActiveTaskID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.activeTaskID
}

// SelectTask выбирает задание (меняет его статус на "в работе")
func (s *TaskService) SelectTask(taskID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если уже есть активное задание с другим ID, запрещаем выбор нового
	if s.activeTaskID != 0 && s.activeTaskID != taskID {
		return fmt.Errorf("уже выполняется задание ID=%d, завершите его перед выбором нового", s.activeTaskID)
	}

	// Если это то же самое задание, что уже активно - ничего не делаем
	if s.activeTaskID == taskID {
		return nil
	}

	// Получаем информацию о задании
	task, err := s.factoryClient.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("ошибка при получении информации о задании: %w", err)
	}

	// Устанавливаем ID активного задания
	s.activeTaskID = taskID

	// Если задание еще не в работе, обновляем его статус через API
	if task.Status != models.TaskStatusInProgress {
		err = s.factoryClient.UpdateTaskStatus(taskID, models.TaskStatusInProgress)
		if err != nil {
			// Логируем ошибку, но продолжаем работу (задание уже считается выбранным)
			return fmt.Errorf("задание выбрано, но возникла ошибка при обновлении статуса: %w", err)
		}
	}

	return nil
}

// FinishTask завершает задание
func (s *TaskService) FinishTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskID := s.activeTaskID
	if taskID == 0 {
		return fmt.Errorf("задание не выбрано")
	}

	// Проверяем, что это активное задание
	if s.activeTaskID != taskID {
		return fmt.Errorf("задание с ID %d не является активным", taskID)
	}

	// Обновляем статус задания на "завершено"
	err := s.factoryClient.UpdateTaskStatus(taskID, models.TaskStatusCompleted)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса задания: %w", err)
	}

	// Сбрасываем активное задание
	s.activeTaskID = 0

	return nil
}
