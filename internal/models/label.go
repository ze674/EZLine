package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// LabelData содержит всю информацию для этикетки
type LabelData struct {
	// Информация о продукте
	Article     string // Артикул
	GTIN        string // GTIN продукта
	Header      string // Шапка этикетки
	Name        string // Название для этикетки
	Standard    string // ТУ/ГОСТ
	Weight      string // Вес единицы (г)
	QuantityBox string // Количество в коробке (шт)
	WeightBox   string // Вес коробки (кг)

	// Информация о задании
	Date        string // Дата производства в формате ДД.ММ.ГГГГ
	BatchNumber string // Номер партии

	// Дополнительная информация
	Packer       string // Упаковщик
	SerialNumber string // Серийный номер

	// Предварительно обработанные данные для шаблона
	BarcodeDate          string // Дата для штрих-кода (ГГММДД)
	FormattedBatchNumber string // Отформатированный номер партии (с ведущими нулями)
	DmData               string // Данные для DataMatrix
	BarcodeText          string // Человекочитаемый текст штрих-кода
	Barcode128Data       string // Данные для Code 128
}

// LabelBuilder предоставляет интерфейс для пошагового построения этикетки
type LabelBuilder struct {
	label LabelData
}

// NewLabelBuilder создает новый экземпляр билдера этикетки
func NewLabelBuilder() *LabelBuilder {
	return &LabelBuilder{
		label: LabelData{},
	}
}

// WithProduct добавляет информацию о продукте из структуры Product
func (b *LabelBuilder) WithProduct(product Product) *LabelBuilder {
	if product.LabelData != "" {
		// Парсим JSON из поля LabelData
		var labelData LabelData
		json.Unmarshal([]byte(product.LabelData), &labelData)

		b.label.Article = labelData.Article
		b.label.GTIN = labelData.GTIN
		b.label.Header = labelData.Header
		b.label.Name = labelData.Name
		b.label.Standard = labelData.Standard
		b.label.Weight = labelData.Weight
		b.label.QuantityBox = labelData.QuantityBox
		b.label.WeightBox = labelData.WeightBox
	}

	return b
}

// WithTask добавляет информацию о задании
func (b *LabelBuilder) WithTask(task Task) *LabelBuilder {
	b.label.Date = task.Date
	b.label.BatchNumber = task.BatchNumber

	// Подготавливаем преобразованные данные
	b.label.BarcodeDate = FormatBarcodeDate(task.Date)
	b.label.FormattedBatchNumber = FormateBatchNumber(task.BatchNumber)
	return b
}

// WithPacker устанавливает имя упаковщика
func (b *LabelBuilder) WithPacker(packer string) *LabelBuilder {
	b.label.Packer = packer
	return b
}

// WithSerialNumber устанавливает серийный номер
func (b *LabelBuilder) WithSerialNumber(sn string) *LabelBuilder {
	b.label.SerialNumber = sn

	// Генерируем комбинированные данные, которые зависят от серийного номера
	if b.label.GTIN != "" && b.label.BarcodeDate != "" && b.label.FormattedBatchNumber != "" {
		b.label.DmData = b.label.GTIN + b.label.BarcodeDate + b.label.FormattedBatchNumber + sn
		b.label.BarcodeText = "(01)" + b.label.GTIN + "(11)" + b.label.BarcodeDate + "(10)" + b.label.FormattedBatchNumber
		b.label.Barcode128Data = "01" + b.label.GTIN + "11" + b.label.BarcodeDate + "10" + b.label.FormattedBatchNumber
	}

	return b
}

// Build создает окончательную структуру этикетки
func (b *LabelBuilder) Build() LabelData {
	return b.label
}

// FormatBarcodeDate преобразует дату из формата ДД.ММ.ГГГГ в формат ГГММДД для штрих-кода
func FormatBarcodeDate(dateStr string) string {
	// Пробуем парсить дату из формата ДД.ММ.ГГГГ
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		// В случае ошибки возвращаем пустую строку или какое-то значение по умолчанию
		return ""
	}

	// Форматируем дату в формат ГГММДД
	return date.Format("060102")
}

func FormateBatchNumber(batchNumberStr string) string {
	// Дополняем нулями слева до 5 цифр
	return fmt.Sprintf("%05s", batchNumberStr)
}

//
