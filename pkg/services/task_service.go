package services

import (
	"fmt"
	"github.com/ze674/EZLine/pkg/api"
	"github.com/ze674/EZLine/pkg/models"
	"sync"
)

type TaskService struct {
	mu            sync.Mutex
	factoryClient *api.FactoryClient
	lineID        int
	activeTask    *models.Task // Указатель на текущее задание
}

// NewTaskService создает новый сервис для управления заданиями
func NewTaskService(factoryClient *api.FactoryClient, lineID int) *TaskService {
	return &TaskService{
		factoryClient: factoryClient,
		lineID:        lineID,
		activeTask:    nil,
	}
}

// GetLineID возвращает ID производственной линии,
// привязанной к этому сервису.
func (s *TaskService) GetLineID() int {
	return s.lineID
}

func (s *TaskService) GetTasks() ([]models.Task, error) {
	return s.factoryClient.GetTasks(s.lineID)
}

// GetActiveTaskID возвращает ID текущего активного задания.
// Если нет активного задания, возвращает 0.
func (s *TaskService) GetActiveTaskID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTask == nil {
		return 0
	}
	return s.activeTask.ID
}

// GetActiveTask возвращает информацию о текущем активном задании.
// Если нет активного задания, возвращает nil.
func (s *TaskService) GetActiveTask() *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если нет активного задания, возвращаем nil
	if s.activeTask == nil {
		return nil
	}

	// Делаем копию задания, чтобы избежать проблем с параллельным доступом
	taskCopy := *s.activeTask
	return &taskCopy
}

// HasActiveTask проверяет, есть ли активное задание.
// Возвращает true, если есть активное задание.
func (s *TaskService) HasActiveTask() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.activeTask != nil
}

// GetTaskByID получает подробную информацию о конкретном задании по ID.
// Метод делает запрос к API для получения актуальных данных.
func (s *TaskService) GetTaskByID(taskID int) (models.Task, error) {
	// Делегируем запрос к API клиенту
	return s.factoryClient.GetTaskByID(taskID)
}

// StartTask устанавливает задание как активное.
// Этот метод:
// 1. Проверяет возможность запуска задания
// 2. Получает свежие данные о задании
// 3. Обновляет внутреннее состояние сервиса
// 4. Обновляет статус задания на сервере
func (s *TaskService) StartTask(taskID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка: если уже есть активное задание с другим ID,
	// запрещаем запуск нового
	if s.activeTask != nil && s.activeTask.ID != taskID {
		return fmt.Errorf("уже выполняется задание ID=%d, завершите его перед началом нового", s.activeTask.ID)
	}

	// Если это то же самое задание, что уже активно - ничего не делаем
	if s.activeTask != nil && s.activeTask.ID == taskID {
		return nil
	}

	// Получаем свежую информацию о задании из API
	task, err := s.factoryClient.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("ошибка при получении информации о задании: %w", err)
	}

	// Устанавливаем задание как активное во внутреннем состоянии
	s.activeTask = &task

	// Если задание еще не в работе, обновляем его статус через API
	if task.Status != models.TaskStatusInProgress {
		err = s.factoryClient.UpdateTaskStatus(taskID, models.TaskStatusInProgress)
		if err != nil {
			// Логируем ошибку, но продолжаем работу (задание уже считается запущенным)
			return fmt.Errorf("задание запущено, но возникла ошибка при обновлении статуса: %w", err)
		}
	}

	return nil
}

// FinishTask завершает текущее активное задание.
// Метод:
// 1. Проверяет наличие активного задания
// 2. Обновляет статус задания на сервере
// 3. Очищает внутреннее состояние активного задания
func (s *TaskService) FinishTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, есть ли активное задание
	if s.activeTask == nil {
		return fmt.Errorf("нет активного задания")
	}

	// Запоминаем ID задания перед сбросом
	taskID := s.activeTask.ID

	// Обновляем статус задания на "завершено" через API
	err := s.factoryClient.UpdateTaskStatus(taskID, models.TaskStatusCompleted)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса задания: %w", err)
	}

	// Сбрасываем информацию об активном задании
	s.activeTask = nil

	return nil
}

// UpdateTaskStatus обновляет статус задания с указанным ID.
// Если это текущее активное задание, также обновляет его статус в памяти.
func (s *TaskService) UpdateTaskStatus(taskID int, status string) error {
	// Сначала обновляем статус на сервере
	err := s.factoryClient.UpdateTaskStatus(taskID, status)
	if err != nil {
		return err
	}

	// Если это текущее активное задание, обновляем его статус в памяти
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTask != nil && s.activeTask.ID == taskID {
		s.activeTask.Status = status
	}

	return nil
}
