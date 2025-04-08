package processors

import (
	"context"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"github.com/ze674/EZLine/internal/repository"
	"github.com/ze674/EZLine/internal/validator"
	"strconv"
	"strings"
	"sync"
)

var serial = 1

type LayerAggregationProcessor struct {
	mu                  sync.Mutex
	running             bool
	cancelFunc          context.CancelFunc
	product             *models.Product
	task                *models.Task
	dataService         DataService
	triggerSource       TriggerSource
	camera              CodeReader
	printer             Printer
	codeValidator       *validator.CodeValidator
	itemRepository      *repository.ItemRepository
	containerRepository *repository.ContainerRepository
}

func NewLayerAggregationProcessor(dataService DataService, scanner CodeReader, source TriggerSource) *LayerAggregationProcessor {
	return &LayerAggregationProcessor{
		dataService:         dataService,
		camera:              scanner,
		triggerSource:       source,
		itemRepository:      repository.NewItemRepository(),
		containerRepository: repository.NewContainerRepository(),
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
	fmt.Println(p.task)

	err = p.connect()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.codeValidator = validator.NewCodeValidator(p.product.GTIN, 31)

	// Создаем контекст, который можно будет отменить при остановке
	ctx, cancel := context.WithCancel(context.Background())
	p.cancelFunc = cancel

	// Запускаем источник
	err = p.triggerSource.WaitSignal(ctx)
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
	if p.printer != nil {

		// Закрываем соединение с принтером
		if err = p.printer.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if p.camera != nil {
		// Закрываем соединение с камерой
		if err = p.camera.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	p.running = false
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
	fmt.Println("Connecting to printer and camera")
	var err error

	//Подключиться к камере
	if p.camera != nil {
		err = p.camera.Connect()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	fmt.Println("Connected to camera")
	//Подключиться к принтеру
	if p.printer != nil {
		err = p.printer.Connect()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	fmt.Println("Connected to printer")
	return nil
}

func (p *LayerAggregationProcessor) runScanningLoop(ctx context.Context) {

	for {
		select {
		case <-p.triggerSource.SignalChan():
			fmt.Println("Scanning started")

			// Сканируем слой
			codes, err := p.scanLayer()
			if err != nil || codes == nil {
				continue
			}

			// Проверяем количество кодов
			//TODO: Сравнить с кол-вом продуктов в коробе
			if len(codes) != 2 {
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

			serialStr := strconv.Itoa(serial)

			err = p.SaveContainerWithItems(serialStr, codes)
			if err != nil {
				continue
			}
			serial++

			//TODO: Проверяем уникальность кодов в задании

			//TODO: Добавляем в список кодов для проверки на уникальность
			//TODO: Печатаем этикетку

		case <-ctx.Done():
			return

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

// SaveContainerWithItems сохраняет контейнер и связанные с ним товары в базу данных
func (p *LayerAggregationProcessor) SaveContainerWithItems(containerCode string, itemCodes []string) error {
	// 1. Создаем контейнер
	containerID, err := p.containerRepository.CreateContainer(
		containerCode,
		p.task.ID,
		repository.StatusCreated,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания контейнера: %w", err)
	}

	// 2. Создаем товары и привязываем их к контейнеру
	for _, code := range itemCodes {
		// Создаем запись о товаре
		itemID, err := p.itemRepository.CreateItem(
			code,
			p.task.ID,
			repository.StatusScanned,
		)
		if err != nil {
			return fmt.Errorf("ошибка создания товара: %w", err)
		}

		// Привязываем товар к контейнеру
		err = p.itemRepository.AssignItemToContainer(itemID, containerID)
		if err != nil {
			return fmt.Errorf("ошибка привязки товара к контейнеру: %w", err)
		}
	}

	return nil
}
