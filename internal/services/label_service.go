package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ze674/EZLine/internal/adapters"
	"github.com/ze674/EZLine/internal/models"
)

// LabelService - сервис для работы с этикетками
type LabelService struct {
	printer      *adapters.Printer
	templatePath string
	Packer       string
	connected    bool
}

// NewLabelService создает новый экземпляр сервиса печати этикеток
func NewLabelService(printer *adapters.Printer, templatePath, defaultPacker string) *LabelService {
	return &LabelService{
		printer:      printer,
		templatePath: templatePath,
		Packer:       defaultPacker,
	}
}

// Connect устанавливает соединение с принтером
func (s *LabelService) Connect() error {
	if s.connected {
		return nil // Уже подключены
	}

	err := s.printer.Connect()
	if err != nil {
		return fmt.Errorf("ошибка при подключении к принтеру: %w", err)
	}

	s.connected = true
	return nil
}

// Close закрывает соединение с принтером
func (s *LabelService) Close() error {
	if !s.connected {
		return nil // Уже отключены
	}

	err := s.printer.Close()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии соединения с принтером: %w", err)
	}

	s.connected = false
	return nil
}

// PrintLabel подготавливает и печатает этикетку на основе данных задания и продукта
func (s *LabelService) PrintLabel(task *models.Task, labelData *models.LabelData, sn string) error {
	// Создаем данные для шаблона
	data := s.prepareTemplateData(task, labelData, sn)

	// Загружаем шаблон
	tmplPath := filepath.Join(s.templatePath, "standard.txt")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке шаблона: %w", err)
	}

	// Заполняем шаблон данными
	result := new(strings.Builder)
	err = tmpl.Execute(result, data)
	if err != nil {
		return fmt.Errorf("ошибка при заполнении шаблона: %w", err)
	}

	// Проверяем, что соединение установлено
	if !s.connected {
		return fmt.Errorf("соединение с принтером не установлено")
	}

	err = s.printer.Send(result.String())
	if err != nil {
		return fmt.Errorf("ошибка при отправке этикетки на печать: %w", err)
	}

	return nil
}

// prepareTemplateData подготавливает данные для шаблона этикетки
func (s *LabelService) prepareTemplateData(task *models.Task, labelData *models.LabelData, sn string) *models.LabelTemplateData {
	// Генерируем серийный номер на основе номера партии и, например, времени
	serialNumber := sn

	// Форматируем дату для штрих-кода
	barcodeDate := models.FormatBarcodeDate(task.Date)
	batchNumber := models.FormateBatchNumber(task.BatchNumber)

	// Заполняем структуру данными
	data := &models.LabelTemplateData{
		// Данные из LabelData
		Article:     labelData.Article,
		GTIN:        labelData.GTIN,
		Header:      labelData.Header,
		Name:        labelData.LabelName,
		Standard:    labelData.Standard,
		Weight:      labelData.UnitWeight,
		QuantityBox: labelData.BoxQuantity,
		WeightBox:   labelData.BoxWeight,

		// Данные из задания
		Date:        task.Date,
		BatchNumber: batchNumber,
		BarcodeDate: barcodeDate,

		// Дополнительные данные
		Packer:       s.Packer,
		SerialNumber: serialNumber,
	}

	return data
}

// ChangePacker изменяет имя упаковщика
func (s *LabelService) ChangePacker(newPacker string) {
	s.Packer = newPacker
}
