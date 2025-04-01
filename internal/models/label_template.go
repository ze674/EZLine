package models

import (
	"fmt"
	"time"
)

// LabelTemplateData содержит все данные для шаблона этикетки
type LabelTemplateData struct {
	// Информация из LabelData
	Article     string // Артикул
	GTIN        string // GTIN продукта
	Header      string // Шапка этикетки
	Name        string // Название для этикетки
	Standard    string // ТУ/ГОСТ
	Weight      string // Вес единицы (г)
	QuantityBox string // Количество в коробке (шт)
	WeightBox   string // Вес коробки (кг)

	// Информация из задания
	Date        string // Дата производства в формате ДД.ММ.ГГГГ
	BatchNumber string // Номер партии
	BarcodeDate string // Дата для штрих-кода (ГГММДД)

	// Дополнительная информация
	Packer       string // Упаковщик
	SerialNumber string // Серийный номер
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
