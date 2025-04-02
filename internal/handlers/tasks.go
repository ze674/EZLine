package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/internal/services"
	"github.com/ze674/EZLine/templates"
	"net/http"
	"strconv"
)

// Добавляем новое поле в структуру TaskHandler
type TaskHandler struct {
	taskService *services.TaskService
	scanService *services.ProcessTaskService // Добавляем сервис сканирования
}

// Обновляем конструктор
func NewTaskHandler(taskService *services.TaskService, scanService *services.ProcessTaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		scanService: scanService,
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

// Обновляем ActiveTaskHandler для передачи статуса сканирования в шаблон
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

	// Проверяем, запущено ли сканирование
	isScanning := h.scanService.IsRunning()
	packer := h.scanService.GetPacker()

	// Отображаем шаблон активного задания
	component := templates.ActiveTask(task, isScanning, packer)

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

// Добавляем обработчик для запуска сканирования
func (h *TaskHandler) StartScanningHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли активное задание
	activeTaskID := h.taskService.GetActiveTaskID()
	if activeTaskID == 0 {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Запускаем сканирование
	err := h.scanService.Start(activeTaskID)
	if err != nil {
		http.Error(w, "Ошибка запуска сканирования: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем обратно на страницу активного задания
	http.Redirect(w, r, "/active-task", http.StatusSeeOther)
}

// Добавляем обработчик для остановки сканирования
func (h *TaskHandler) StopScanningHandler(w http.ResponseWriter, r *http.Request) {
	// Останавливаем сканирование
	h.scanService.Stop()

	// Перенаправляем обратно на страницу активного задания
	http.Redirect(w, r, "/active-task", http.StatusSeeOther)
}

// ChangePacker обрабатывает запрос на изменение упаковщика
func (h *TaskHandler) ChangePackerHandler(w http.ResponseWriter, r *http.Request) {
	// Обрабатываем только POST запросы
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем новое имя упаковщика из формы
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	newPacker := r.FormValue("packer")
	if newPacker == "" {
		http.Error(w, "Имя упаковщика не может быть пустым", http.StatusBadRequest)
		return
	}

	// Меняем упаковщика в сервисе
	h.scanService.ChangePacker(newPacker)

	// Перенаправляем на страницу активного задания
	http.Redirect(w, r, "/active-task", http.StatusSeeOther)
}
