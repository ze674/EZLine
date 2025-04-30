package processors

import (
	"context"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"sync"
)

var sensorChan <-chan struct{}

type Scanner interface {
	Close() error
	Connect() error
	Scan() (string, error)
}

// PLC интерфейс описывает работу с внешним ПЛК
type PLC interface {
	Connect() error
	Close() error
	HandleProductSignal(ctx context.Context) (<-chan struct{}, error)
	RejectorOn() error
	RejectorOff() error
}

type AutomaticSerializationProcessor struct {
	mu          sync.Mutex
	wg          sync.WaitGroup
	cancelFunc  context.CancelFunc
	running     bool
	dataService DataService
	task        *models.Task
	product     *models.Product
	scanner     Scanner
	plc         PLC
}

func NewAutomaticSerializationProcessor(dataService DataService, scanner Scanner, plc PLC) *AutomaticSerializationProcessor {
	return &AutomaticSerializationProcessor{
		plc:         plc,
		scanner:     scanner,
		dataService: dataService,
	}
}

func (p *AutomaticSerializationProcessor) Start(TaskID int) error {
	op := "processors.AutomaticSerializationProcessor.Start"
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

	fmt.Println(p.task)
	fmt.Println(p.product)

	// Создаем контекст, который можно будет отменить при остановке
	ctx, cancel := context.WithCancel(context.Background())
	p.cancelFunc = cancel

	if err := p.connect(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sensorChan, err = p.plc.HandleProductSignal(ctx)
	p.wg.Add(1)
	go p.run(ctx)
	p.running = true

	return nil
}

func (p *AutomaticSerializationProcessor) Stop() error {
	op := "processors.AutomaticSerializationProcessor.Stop"

	//var err error

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

	p.wg.Wait() // Ожидаем завершения работы горутины

	if err := p.disconnect(); err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	p.running = false
	return nil

}

func (p *AutomaticSerializationProcessor) IsRunning() bool {
	return p.running
}

func (p *AutomaticSerializationProcessor) getData(TaskID int) error {
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

func (p *AutomaticSerializationProcessor) connect() error {
	op := "processors.AutomaticSerializationProcessor.connect"

	var err error

	if p.scanner != nil {
		if err = p.scanner.Connect(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if p.plc != nil {
		if err = p.plc.Connect(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (p *AutomaticSerializationProcessor) disconnect() error {
	op := "processors.AutomaticSerializationProcessor.disconnect"

	var err error

	if p.scanner != nil {
		if err = p.scanner.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if p.plc != nil {
		if err = p.plc.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (p *AutomaticSerializationProcessor) run(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			<-sensorChan
			res, err := p.scanner.Scan()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(res)
		}
	}
}
