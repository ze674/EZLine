package models

// LabelTemplateData содержит все данные для шаблона этикетки
type LabelTemplateData struct {
	productData  LabelProductData // Информация о продукте
	taskData     LabelTaskData    // Информация о задании
	Packer       string
	SerialNumber string // Серийный номер

	//Трансформированные данные готовые для отображения
	ContainerBarcode string // Аггрегационный код
	Barcode128Data   string // Стандартный штрих код
	Barcode128Text   string // Человекочитаемая часть
}

type LabelProductData struct {
	Article     string // Артикул
	Header      string // Шапка этикетки
	Name        string // Название для этикетки
	Standard    string // ТУ/ГОСТ
	Weight      string // Вес единицы (г)
	QuantityBox string // Количество в коробке (шт)
	WeightBox   string // Вес коробки (кг)
	GTIN        string // GTIN короба
}

type LabelTaskData struct {
	Date        string // Дата производства в формате ДД.ММ.ГГГГ
	BatchNumber string // Номер партии
}
