package processors

import (
	"context"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"github.com/ze674/EZLine/internal/validator"
	"strings"
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
	codeValidator *validator.CodeValidator
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

	p.codeValidator = validator.NewCodeValidator(p.product.GTIN, 31)

	// Создаем контекст, который можно будет отменить при остановке
	ctx, cancel := context.WithCancel(context.Background())
	p.cancelFunc = cancel

	// Запускаем источник
	err = p.triggerSource.Start(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

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

	for {
		select {
		case <-p.triggerSource.Signal():

			// Сканируем слой
			codes, err := p.scanLayer()
			if err != nil || codes == nil {
				continue
			}

			// Проверяем количество кодов
			//TODO: Сравнить с кол-вом продуктов в коробе
			if len(codes) != 6 {
				continue
			}

			// Проверяем наличие дубликатов в слое
			if p.checkDuplicatesInLayer(codes) {
				continue
			}

			//Валидируем коды
			validationResult := p.codeValidator.ValidateCodes(codes)

			if !validationResult.Valid {
				continue
			}

			//TODO: Проверяем уникальность кодов в задании

			//TODO: Сохраняем агрегат в базу
			//TODO: Добавляем в список кодов для проверки на уникальность
			//TODO: Печатаем этикетку

		case <-ctx.Done():

		}
	}
}

func (p *LayerAggregationProcessor) scanLayer() ([]string, error) {
	op := "processors.LayerAggregationProcessor.scanLayer"

	resp, err := p.camera.Scan()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp == "NoRead" {
		return nil, nil
	}

	codes := strings.Fields(resp)

	return codes, nil
}

// Проверяем наличие дубликатов в слое
func (p *LayerAggregationProcessor) checkDuplicatesInLayer(codes []string) bool {
	seen := make(map[string]bool)

	for _, code := range codes {
		if seen[code] {
			// Нашли дубликат, можно сразу вернуть true
			return true
		}
		seen[code] = true
	}

	// Дубликатов не найдено
	return false
}
