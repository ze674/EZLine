package models

// ScanStats содержит статистику сканирования
type ScanStats struct {
	TotalCodes    int      // Общее количество кодов
	ValidCodes    int      // Количество валидных кодов
	InvalidCodes  int      // Количество невалидных кодов
	RecentCodes   []string // Последние коды (до 10)
	RecentResults []bool   // Результаты валидации последних кодов
	LastCode      string   // Последний отсканированный код
	LastResult    bool     // Результат валидации последнего кода
	LastError     string   // Ошибка валидации последнего кода
}
