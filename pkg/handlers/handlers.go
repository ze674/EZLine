package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/pkg/api"
	"github.com/ze674/EZLine/templates"
	"net/http"
)

// Глобальные настройки
var (
	FactoryClient *api.FactoryClient
	LineID        int // ID производственной линии из конфигурации
	ActiveTaskID  int // 0 - если нет текущего задания
)

// Init инициализирует обработчики
func Init(factoryURL string, lineID int) {
	FactoryClient = api.NewFactoryClient(factoryURL)
	LineID = lineID
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
		r.Get("/{id}", StartScanningHandler)     // начало сканирования
		r.Post("/finish", FinishScanningHandler) // завершение сканирования
	})
}

// homeHandler отображает домашнюю страницу
func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.Home()

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}
