// internal/services/serial_generator.go
package services

import (
	"fmt"
	"github.com/ze674/EZLine/internal/repository"
	"sync"
)

// SerialGenerator отвечает за генерацию серийных номеров для контейнеров
type SerialGenerator struct {
	mu                  sync.Mutex
	lastSerial          int  // последний сгенерированный серийный номер
	taskID              int  // ID текущего задания
	initialized         bool // флаг инициализации
	containerRepository *repository.ContainerRepository
}

// NewSerialGenerator создает новый генератор серийных номеров
func NewSerialGenerator() *SerialGenerator {
	return &SerialGenerator{
		containerRepository: repository.NewContainerRepository(),
		lastSerial:          0,
		initialized:         false,
	}
}

// Initialize инициализирует генератор для конкретного задания
func (g *SerialGenerator) Initialize(taskID int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.taskID = taskID

	// Получаем последний серийный номер для задания
	lastSerial, err := g.containerRepository.GetLastSerialNumber(taskID)
	if err != nil {
		return fmt.Errorf("ошибка при получении последнего серийного номера: %w", err)
	}

	g.lastSerial = lastSerial
	g.initialized = true
	return nil
}

// GenerateSerial генерирует новый серийный номер и возвращает его в строковом формате
func (g *SerialGenerator) GenerateSerial() (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.initialized {
		return 0, fmt.Errorf("генератор не инициализирован")
	}

	// Увеличиваем счетчик
	g.lastSerial++

	return g.lastSerial, nil
}

// GetLastSerial возвращает последний сгенерированный серийный номер
func (g *SerialGenerator) GetLastSerial() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.lastSerial
}
