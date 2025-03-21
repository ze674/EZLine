package config

import (
	"encoding/json"
	"os"
)

// Config содержит настройки приложения
type Config struct {
	FactoryURL string `json:"factory_url"` // URL API EZFactory
	LineID     int    `json:"line_id"`     // ID производственной линии
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		FactoryURL: "http://localhost:8081",
		LineID:     1,
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
