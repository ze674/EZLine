package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/templates"
	"net/http"
	"strconv"
)

// StartScanningHandler обрабатывает запрос на начало сканирования для выбранного задания
func StartScanningHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID задания", http.StatusBadRequest)
		return
	}

	// Если уже есть активное задание, не даем начать новое
	if ActiveTaskID != 0 && ActiveTaskID != taskID {
		http.Error(w, "Уже выполняется другое задание. Завершите его перед началом нового.", http.StatusBadRequest)
		return
	}

	// Получаем информацию о задании
	task, err := FactoryClient.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем задание как активное
	ActiveTaskID = taskID

	// Обновляем статус задания на "в работе", если оно еще не в работе
	if task.Status != "в работе" {
		err = FactoryClient.UpdateTaskStatus(taskID, "в работе")
		if err != nil {
			http.Error(w, "Ошибка при обновлении статуса задания: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Отображаем экран сканирования
	component := templates.ScanningScreen(task)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// FinishScanningHandler обрабатывает запрос на завершение сканирования
func FinishScanningHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что есть активное задание
	if ActiveTaskID == 0 {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Обновляем статус задания на "завершено"
	err := FactoryClient.UpdateTaskStatus(ActiveTaskID, "завершено")
	if err != nil {
		http.Error(w, "Ошибка при обновлении статуса задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сбрасываем активное задание
	ActiveTaskID = 0

	// Перенаправляем на список заданий
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
