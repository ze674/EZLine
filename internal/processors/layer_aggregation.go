package processors

import (
	"context"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"sync"
)

type LayerAggregationProcessor struct {
	mu            sync.Mutex
	running       bool
	cancelFunc    context.CancelFunc
	product       *models.Product
	task          *models.Task
	dataService   DataService
	triggerSource TriggerSource
	camera        CodeReader
	printer       Printer
}

func NewLayerAggregationProcessor(dataService DataService, scanner CodeReader, source TriggerSource) *LayerAggregationProcessor {
	return &LayerAggregationProcessor{
		dataService:   dataService,
		camera:        scanner,
		triggerSource: source,
	}
}

func (p *LayerAggregationProcessor) Start(TaskID int) error {
	op := "processors.LayerAggregationProcessor.Start"

	var err error

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	err = p.getData(TaskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = p.connect()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Создаем контекст, который можно будет отменить при остановке
	ctx, cancel := context.WithCancel(context.Background())
	p.cancelFunc = cancel

	go p.runScanningLoop(ctx)

	p.running = true

	return nil
}

// Stop останавливает процессор
func (p *LayerAggregationProcessor) Stop() error {
	op := "processors.LayerAggregationProcessor.Stop"

	var err error

	if !p.running {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Отменяем контекст, что приведет к завершению цикла сканирования
	if p.cancelFunc != nil {
		p.cancelFunc()
		p.cancelFunc = nil
	}

	// Закрываем соединение с принтером
	if err = p.printer.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// Закрываем соединение с камерой
	if err = p.camera.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *LayerAggregationProcessor) IsRunning() bool {
	return p.running
}

func (p *LayerAggregationProcessor) getData(TaskID int) error {
	op := "processors.LayerAggregationProcessor.init"

	var err error

	//Загрузить информацию о задании.
	task, err := p.dataService.GetTaskByID(TaskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.task = &task

	//Загрузить информацию о продукте.
	product, err := p.dataService.GetProductByID(p.task.ProductID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.product = &product

	return nil
}

func (p *LayerAggregationProcessor) connect() error {
	op := "processors.LayerAggregationProcessor.connect"

	var err error

	p.mu.Lock()
	defer p.mu.Unlock()

	//Подключиться к камере
	err = p.camera.Connect()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//Подключиться к принтеру
	err = p.printer.Connect()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *LayerAggregationProcessor) runScanningLoop(ctx context.Context) {
	scanTrigger, err := p.triggerSource.Start(ctx)
	if err != nil {
		go p.Stop()
	}
	for {
		select {
		case <-scanTrigger:

		case <-ctx.Done():

		}
	}
}
