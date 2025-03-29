package services

import (
	"context"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"sync"
	"time"
)

// services/data_access.go
type DataService interface {
	GetTaskByID(taskID int) (models.Task, error)
	GetProductByID(productID int) (models.Product, error)
	// Можно добавить другие методы по мере необходимости
}

type scanner interface {
	Scan() (string, error)
	Connect() error
	Close() error
}

type ProcessTaskService struct {
	mu             sync.Mutex
	interval       time.Duration
	scanner        scanner
	DataService    DataService
	running        bool
	cancelFunc     context.CancelFunc
	currentTask    *models.Task
	currentProduct *models.Product
}

func NewProcessTaskService(dataService DataService, scanner scanner, interval time.Duration) *ProcessTaskService {
	return &ProcessTaskService{
		DataService: dataService,
		scanner:     scanner,
		interval:    interval,
		running:     false,
		currentTask: nil,
	}
}

// Start запускает сервис сканирования
func (s *ProcessTaskService) Start(id int) error {
	var err error
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если сервис уже запущен, ничего не делаем
	if s.running {
		return nil
	}
	fmt.Printf("Start task ID=%d\n", id)
	task, err := s.DataService.GetTaskByID(id)
	if err != nil {
		fmt.Printf("Error get task ID=%d: %s\n", id, err.Error())
		return err
	}

	s.currentTask = &task

	fmt.Errorf("Start product ID=%d\n", s.currentTask.ProductID)
	product, err := s.DataService.GetProductByID(s.currentTask.ProductID)
	if err != nil {
		fmt.Printf("Error get product ID=%d: %s\n", s.currentTask.ProductID, err.Error())
		return err
	}

	s.currentProduct = &product

	// Создаем контекст, который можно будет отменить при остановке
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	// Запускаем цикл сканирования в отдельной горутине
	go s.runScanLoop(ctx)

	s.running = true
	return nil
}

// Stop останавливает сервис сканирования
func (s *ProcessTaskService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	// Отменяем контекст, что приведет к завершению цикла сканирования
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}

	s.running = false
}

// runScanLoop запускает цикл сканирования в фоновом режиме
func (s *ProcessTaskService) runScanLoop(ctx context.Context) {
	// Подключаемся к сканеру
	if err := s.scanner.Connect(); err != nil {
		// В MVP просто отмечаем, что сервис не запущен
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return
	}

	// Не забываем закрыть соединение при завершении
	defer s.scanner.Close()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Выполняем сканирование
			result, err := s.scanner.Scan()
			if err != nil {
				// В MVP просто продолжаем работу
				continue
			}

			// Используем информацию о задании и продукте
			fmt.Printf("Задание: %d, Продукт: %s, Результат сканирования: %s\n",
				s.currentTask.ID, s.currentProduct.Name, result)

		case <-ctx.Done():
			// Контекст был отменен, завершаем работу
			return
		}
	}
}

// GetCurrentTask возвращает информацию о текущем задании
func (s *ProcessTaskService) GetCurrentTask() *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.currentTask
}

// IsRunning возвращает текущее состояние сервиса
func (s *ProcessTaskService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// GetCurrentProduct возвращает информацию о текущем продукте
