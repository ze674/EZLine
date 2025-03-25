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
	TasksHandlers *TaskHandler
)

// Init инициализирует обработчики
func Init(tasksHandlers *TaskHandler) {
	TasksHandlers = tasksHandlers
}

// SetupRoutes настраивает маршруты HTTP
func SetupRoutes(r chi.Router) {
	// Домашняя страница
	r.Get("/", homeHandler)

	// Маршруты для работы с заданиями
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", TasksHandlers.ListTasksHandler) // список заданий
	})

	// Маршруты для сканирования
	r.Route("/scanning", func(r chi.Router) {
		r.Get("/{id}", TasksHandlers.StartTaskHandler)             // запуск задания
		r.Get("/view", TasksHandlers.ViewScanningHandler)          // представление сканирования
		r.Post("/start-roll", TasksHandlers.StartRollHandler)      // начало сканирования ролика
		r.Post("/finish-roll", TasksHandlers.FinishRollHandler)    // завершение сканирования ролика
		r.Get("/refresh-stats", TasksHandlers.RefreshStatsHandler) // обновление статистики
		r.Post("/finish", TasksHandlers.FinishTaskHandler)         // завершение текущее задание
	})
}

// homeHandler отображает домашнюю страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
