// cmd/ezline/main.go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ze674/EZLine/pkg/api"
	"github.com/ze674/EZLine/pkg/config"
	"github.com/ze674/EZLine/pkg/handlers"
	"github.com/ze674/EZLine/pkg/services"
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

	factoryClient := api.NewFactoryClient(cfg.FactoryURL)
	taskService := services.NewTaskService(factoryClient, cfg.LineID)
	scanService := services.NewScanService(cfg.ScannerAddress, cfg.StoragePath, cfg.CodeLength)
	taskHandlers := handlers.NewTaskHandler(taskService, scanService)
	// Инициализируем обработчики
	handlers.Init(taskHandlers)

	// Создаем роутер
	r := chi.NewRouter()

	// Middleware
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Настраиваем маршруты
	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	handlers.SetupRoutes(r)

	// Запускаем сервер
	log.Printf("Запуск сервера на http://localhost:8080 (Линия ID: %d, EZFactory: %s)",
		cfg.LineID, cfg.FactoryURL)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
