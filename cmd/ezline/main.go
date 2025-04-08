// cmd/ezline/main.go
package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ze674/EZLine/internal/adapters"
	"github.com/ze674/EZLine/internal/api"
	"github.com/ze674/EZLine/internal/config"
	"github.com/ze674/EZLine/internal/database"
	"github.com/ze674/EZLine/internal/handlers"
	"github.com/ze674/EZLine/internal/processors"
	"github.com/ze674/EZLine/internal/services"
	"github.com/ze674/EZLine/internal/utils"
	"log"
	"net/http"
	"time"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v. Используем значения по умолчанию.", err)
		cfg = config.DefaultConfig()
	}

	fmt.Printf("Scan coomand is %s, Noread answer is %s", cfg.ScanCommand, cfg.AnswerNoRead)

	// Инициализируем TCP сканер
	scanner := adapters.NewScanner(cfg.ScannerAddress, cfg.ScanCommand) // Пустая команда или команда по вашему выбору
	//printer := adapters.NewPrinter(cfg.PrinterAddress)
	// Инициализируем соединение с базой данных

	if err := database.Connect(cfg.DbPath); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()
	// Замените на реальные значения или используйте из конфига

	factoryClient := api.NewFactoryClient(cfg.FactoryURL)

	//// Инициализируем сервис печати этикеток
	//labelService := services.NewLabelService(
	//	printer,
	//	cfg.TemplatePath, // Путь к шаблонам этикеток
	//	"",               // Значение по умолчанию для упаковщика
	//)
	taskService := services.NewTaskService(factoryClient, cfg.LineID)
	//scanService := services.NewProcessTaskService(taskService, labelService, scanner, 2*time.Second)
	//// Передаем сервис сканирования в обработчик заданий
	//

	triggerService := utils.NewTimerTrigger(1000 * time.Millisecond)
	scanService := processors.NewLayerAggregationProcessor(taskService, scanner, triggerService)
	taskHandlers := handlers.NewTaskHandler(taskService, scanService)

	// Создаем роутер
	r := chi.NewRouter()

	// Middleware
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Настраиваем маршруты
	// Обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	handlers.SetupRoutes(r, taskHandlers)

	// Запускаем сервер
	log.Printf("Запуск сервера на http://localhost:8080 (Линия ID: %d, EZFactory: %s)",
		cfg.LineID, cfg.FactoryURL)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
