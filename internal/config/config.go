package config

import (
	"encoding/json"
	"os"
)

// Config содержит настройки приложения
type Config struct {
	FactoryURL     string `json:"factory_url"`           // URL API EZFactory
	LineID         int    `json:"line_id"`               // ID производственной линии
	ScannerAddress string `json:"scanner_address"`       // IP адрес сканнера
	PrinterAddress string `json:"printer_address"`       // IP адрес принтера
	StoragePath    string `json:"storage_path"`          // Путь к файлу хранения файлов
	TemplatePath   string `json:"template_path"`         // Путь к шаблонам этикеток
	DbPath         string `json:"db_path"`               // Путь к базе данных
	CodeLength     int    `json:"code_length"`           // Длина кода
	ScanCommand    string `json:"scanner_answer_noread"` // Команда сканирования
	AnswerNoRead   string `json:"scanner_scan_command"`  // Ответ на команду сканирования
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		FactoryURL:     "http://localhost:8081",
		LineID:         1,
		StoragePath:    "./data",
		ScannerAddress: "127.0.0.1:2001",
	}
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(path string) (Config, error) {
	// Если файл не найден, используем конфигурацию по умолчанию
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultConfig := DefaultConfig()

		// Создаем файл с конфигурацией по умолчанию
		if err := SaveConfig(path, defaultConfig); err != nil {
			return defaultConfig, err
		}

		return defaultConfig, nil
	}

	// Читаем файл
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), err
	}

	// Разбираем JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}

// SaveConfig сохраняет конфигурацию в файл
func SaveConfig(path string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
