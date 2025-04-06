package processors

import (
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"sync"
)

type LayerAggregationProcessor struct {
	mu          sync.Mutex
	running     bool
	product     *models.Product
	task        *models.Task
	dataService DataService
	camera      CodeReader
	printer     Printer
}

func NewLayerAggregationProcessor(dataService DataService, scanner CodeReader) *LayerAggregationProcessor {
	return &LayerAggregationProcessor{
		dataService: dataService,
		camera:      scanner,
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

	p.running = true

	return nil
}

func (p *LayerAggregationProcessor) Stop() error {
	op := "processors.LayerAggregationProcessor.Stop"

	var err error

	p.mu.Lock()
	defer p.mu.Unlock()

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
