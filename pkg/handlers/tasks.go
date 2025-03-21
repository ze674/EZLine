package handlers

import (
	"github.com/ze674/EZLine/templates"
	"net/http"
)

// tasksListHandler отображает список заданий для линии
func tasksListHandler(w http.ResponseWriter, r *http.Request) {
	// Используем ID линии из конфигурации
	lineID := LineID

	// Получаем задания из API EZFactory
	tasks, err := FactoryClient.GetTasks(lineID)
	if err != nil {
		http.Error(w, "Ошибка при получении списка заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим шаблон
	component := templates.TasksList(tasks, lineID)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}
