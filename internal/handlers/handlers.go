// internal/handlers/handlers.go
package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// internal/handlers/handlers.go
func SetupRoutes(r chi.Router, taskHandler *TaskHandler) {
	r.Get("/", homeHandler)

	// Маршруты для заданий
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.ListTasksHandler)              // список заданий
		r.Post("/{id}/select", taskHandler.SelectTaskHandler) // выбор задания
		r.Post("/finish", taskHandler.FinishTaskHandler)      // завершение задания
	})

	// Страница активного задания
	r.Get("/active-task", taskHandler.ActiveTaskHandler)

	// Добавляем маршруты для управления сканированием
	r.Post("/scanning/start", taskHandler.StartScanningHandler)
	r.Post("/scanning/stop", taskHandler.StopScanningHandler)
	//r.Post("/packer/change", taskHandler.ChangePackerHandler)
}

// homeHandler отображает домашнюю страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
