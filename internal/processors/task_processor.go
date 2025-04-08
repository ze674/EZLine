package processors

import (
	"context"
	"github.com/ze674/EZLine/internal/models"
)

// TaskProcessor определяет общий интерфейс для обработки заданий
type TaskProcessor interface {
	// Start запускает процессор
	Start(TaskID int) error

	// Stop останавливает процессор
	Stop() error

	// IsRunning возвращает состояние процессора
	IsRunning() bool
}

type DataService interface {
	GetTaskByID(taskID int) (models.Task, error)
	GetProductByID(productID int) (models.Product, error)
}

type CodeReader interface {
	Scan() (string, error)
	Connect() error
	Close() error
}

type Printer interface {
	Print(data string) error
	Connect() error
	Close() error
}

// TriggerSource представляет источник триггеров для сканирования
type TriggerSource interface {
	Signal() <-chan struct{}
	// Start запускает источник и возвращает канал с сигналами для сканирования
	Start(ctx context.Context) error
	// Stop останавливает источник триггеров
	Stop() error
}
