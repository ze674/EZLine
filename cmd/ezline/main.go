package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ze674/EZLine/pkg/config"
	"github.com/ze674/EZLine/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v. Используем значения по умолчанию.", err)
		cfg = config.DefaultConfig()
	}

	// Инициализируем обработчики с настройками из конфигурации
	handlers.Init(cfg.FactoryURL, cfg.LineID)

	// Создаем роутер
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Настраиваем маршруты
	handlers.SetupRoutes(r)

	// Запускаем сервер
	log.Printf("Запуск сервера на http://localhost:8080 (Линия ID: %d, EZFactory: %s)",
		cfg.LineID, cfg.FactoryURL)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
