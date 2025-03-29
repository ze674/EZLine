package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"github.com/ze674/EZLine/internal/validator"
	"strconv"
	"strings"
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
	codeValidator  *validator.CodeValidator
	running        bool
	cancelFunc     context.CancelFunc
	currentTask    *models.Task
	currentProduct *models.Product
	labelData      *models.LabelData
	QuantityPerBox int
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
	task, err := s.DataService.GetTaskByID(id)
	if err != nil {
		return err
	}

	s.currentTask = &task

	product, err := s.DataService.GetProductByID(s.currentTask.ProductID)
	if err != nil {
		return err
	}

	s.currentProduct = &product

	// Парсим LabelData
	if s.currentProduct.LabelData != "" {
		var labelData models.LabelData
		if err := json.Unmarshal([]byte(s.currentProduct.LabelData), &labelData); err != nil {
			// Только логируем ошибку, но продолжаем работу
			fmt.Printf("Ошибка при парсинге LabelData: %v\n", err)
		} else {
			s.labelData = &labelData
		}
	}

	// Создаем валидатор с параметрами из загруженного продукта
	s.codeValidator = validator.NewCodeValidator(
		s.currentProduct.GTIN,
		31, // Длина кода как константа или из конфигурации
	)

	// Получаем ожидаемое количество кодов из LabelData
	s.QuantityPerBox = 1 // По умолчанию
	if s.labelData != nil && s.labelData.BoxQuantity != "" {
		if qty, err := strconv.Atoi(s.labelData.BoxQuantity); err == nil {
			s.QuantityPerBox = qty
		}
	}

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

func (s *ProcessTaskService) runScanLoop(ctx context.Context) {
	// Подключаемся к сканеру
	if err := s.scanner.Connect(); err != nil {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		fmt.Printf("Ошибка подключения к сканеру: %v\n", err)
		return
	}
	defer s.scanner.Close()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Выполняем сканирование
			response, err := s.scanner.Scan()
			if err != nil {
				fmt.Printf("Ошибка сканирования: %v\n", err)
				continue
			}

			// Если это специальное значение "NoRead", пропускаем
			if response == "NoRead" {
				fmt.Println("Код не прочитан")
				continue
			}

			// Разделяем строку ответа на отдельные коды
			codes := strings.Fields(response) // Разделяем по пробелам

			// Проверяем количество кодов
			if len(codes) != s.QuantityPerBox {
				fmt.Printf("Ошибка: получено %d кодов, ожидалось %d\n", len(codes), s.QuantityPerBox)
				continue
			}

			// Проверяем каждый код с помощью валидатора
			allValid := true
			var invalidCodes []string

			for _, code := range codes {
				if s.codeValidator != nil {
					validationResult := s.codeValidator.ValidateCode(code)
					if !validationResult.Valid {
						allValid = false
						invalidCodes = append(invalidCodes,
							fmt.Sprintf("%s (причина: %s)", code, validationResult.Message))
					}
				}
			}

			// Выводим результаты
			if allValid {
				fmt.Printf("Все %d кодов валидны!\n", len(codes))
				// Здесь можно добавить логику успешной обработки коробки
			} else {
				fmt.Printf("Найдены невалидные коды: %s\n", strings.Join(invalidCodes, ", "))
				// Здесь можно добавить логику отбраковки
			}

		case <-ctx.Done():
			// Контекст отменен, завершаем работу
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
