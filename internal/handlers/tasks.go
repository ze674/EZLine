package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/internal/services"
	"github.com/ze674/EZLine/templates"
	"net/http"
	"strconv"
)

// TaskHandler обрабатывает запросы, связанные с заданиями
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler создает новый обработчик заданий
func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// ListTasksHandler отображает список заданий для линии
func (h *TaskHandler) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем задания из сервиса
	tasks, err := h.taskService.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка при получении списка заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем ID линии и активного задания
	lineID := h.taskService.GetLineID()
	activeTaskID := h.taskService.GetActiveTaskID()

	// Рендерим шаблон
	component := templates.TasksList(tasks, lineID, activeTaskID)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// handlers/task_handler.go
// ActiveTaskHandler показывает страницу с активным заданием
func (h *TaskHandler) ActiveTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли активное задание
	activeTaskID := h.taskService.GetActiveTaskID()
	if activeTaskID == 0 {
		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	// Получаем информацию о задании
	task, err := h.taskService.GetTaskByID(activeTaskID)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Отображаем шаблон активного задания
	component := templates.ActiveTask(task)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// SelectTaskHandler выбирает задание
func (h *TaskHandler) SelectTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID задания", http.StatusBadRequest)
		return
	}

	// Выбираем задание
	err = h.taskService.SelectTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Перенаправляем на страницу активного задания
	http.Redirect(w, r, "/active-task", http.StatusSeeOther)
}

// FinishTaskHandler завершает текущее выбранное задание
func (h *TaskHandler) FinishTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Завершаем текущее активное задание
	err := h.taskService.FinishTask()
	if err != nil {
		http.Error(w, "Ошибка при завершении задания: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на список заданий
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
