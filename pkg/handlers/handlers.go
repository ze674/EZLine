// pkg/handlers/handlers.go
package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/pkg/api"
	"net/http"
)

// Глобальные настройки
var (
	FactoryClient *api.FactoryClient
	LineID        int
	ActiveTaskID  int
)

// Init инициализирует обработчики
func Init(factoryURL string, lineID int, scannerAddress string, storagePath string, codeLength int) {
	FactoryClient = api.NewFactoryClient(factoryURL)
	LineID = lineID

	// Инициализируем модуль сканирования
	InitScanning(scannerAddress, storagePath, codeLength)
}

// SetupRoutes настраивает маршруты HTTP
func SetupRoutes(r chi.Router) {
	// Домашняя страница
	r.Get("/", homeHandler)

	// Маршруты для работы с заданиями
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", tasksListHandler) // список заданий
	})

	// Маршруты для сканирования
	r.Route("/scanning", func(r chi.Router) {
		r.Get("/{id}", StartTaskHandler)
		r.Get("/view", ViewScanningHandler)
		r.Post("/start-roll", StartRollHandler)      // начало сканирования ролика
		r.Post("/finish-roll", FinishRollHandler)    // завершение сканирования ролика
		r.Get("/refresh-stats", RefreshStatsHandler) // обновление статистики
		r.Post("/finish", FinishScanningHandler)     // завершение сканирования
	})
}

// homeHandler отображает домашнюю страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
