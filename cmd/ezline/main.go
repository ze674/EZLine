// cmd/ezline/main.go
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ze674/EZLine/internal/adapters"
	"github.com/ze674/EZLine/internal/api"
	"github.com/ze674/EZLine/internal/config"
	"github.com/ze674/EZLine/internal/database"
	"github.com/ze674/EZLine/internal/handlers"
	"github.com/ze674/EZLine/internal/processors"
	"github.com/ze674/EZLine/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v. Используем значения по умолчанию.", err)
		cfg = config.DefaultConfig()
	}

	// Подключаемся к базе данных
	if err := database.Connect(cfg.DbPath); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	factoryClient := api.NewFactoryClient(cfg.FactoryURL)

	taskService := services.NewTaskService(factoryClient, cfg.LineID)

	camera := adapters.NewScanner(cfg.ScannerAddress, cfg.ScanCommand)
	pusherReg, err := strconv.Atoi(cfg.PusherRegister)
	if err != nil {
		log.Fatal(err)
	}
	sensorReg, err := strconv.Atoi(cfg.SensorRegister)
	if err != nil {
		log.Fatal(err)
	}

	plc := adapters.NewModbusPLC(cfg.PlcAddress, 5*time.Second, 5*time.Millisecond, uint16(sensorReg), uint16(pusherReg), 5)

	scanService := processors.NewAutomaticSerializationProcessor(taskService, camera, plc)

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
