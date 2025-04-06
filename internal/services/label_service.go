// internal/services/label_service.go

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
		connected:    false,
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

// RenderTemplate рендерит этикетку в строку для печати
func (s *LabelService) RenderTemplate(labelData models.LabelData, templateName string) (string, error) {
	if templateName == "" {
		templateName = "standard.txt" // Шаблон по умолчанию
	}

	// Загружаем шаблон
	tmplPath := filepath.Join(s.templatePath, templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("ошибка при загрузке шаблона %s: %w", tmplPath, err)
	}

	// Заполняем шаблон данными
	result := new(strings.Builder)
	err = tmpl.Execute(result, labelData)
	if err != nil {
		return "", fmt.Errorf("ошибка при заполнении шаблона: %w", err)
	}

	return result.String(), nil
}

// Print отправляет подготовленный контент на печать
func (s *LabelService) Print(content string) error {
	// Проверяем, что соединение установлено
	if !s.connected {
		if err := s.Connect(); err != nil {
			return fmt.Errorf("не удалось установить соединение с принтером: %w", err)
		}
	}

	// Отправляем на печать
	err := s.printer.Send(content)
	if err != nil {
		return fmt.Errorf("ошибка при отправке на печать: %w", err)
	}

	return nil
}

// RenderAndPrint комбинирует рендеринг и печать в один удобный метод
func (s *LabelService) RenderAndPrint(labelData models.LabelData, templateName string) error {
	content, err := s.RenderTemplate(labelData, templateName)
	if err != nil {
		return err
	}

	return s.Print(content)
}

// ChangePacker изменяет имя упаковщика
func (s *LabelService) ChangePacker(newPacker string) {
	s.Packer = newPacker
}

// GetPacker возвращает текущее имя упаковщика
func (s *LabelService) GetPacker() string {
	return s.Packer
}

// PrintLabel - удобный метод для печати с использованием Builder
func (s *LabelService) PrintLabel(task *models.Task, product *models.Product, serialNumber string) error {
	// Создаем билдер
	labelBuilder := models.NewLabelBuilder()

	// Добавляем данные
	labelBuilder.WithProduct(*product)
	labelBuilder.WithTask(*task)
	labelBuilder.WithPacker(s.GetPacker())
	labelBuilder.WithSerialNumber(serialNumber)

	// Собираем этикетку
	labelData := labelBuilder.Build()

	// Рендерим и печатаем
	return s.RenderAndPrint(labelData, "standard.txt")
}
