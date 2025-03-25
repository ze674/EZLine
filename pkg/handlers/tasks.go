package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/pkg/models"
	"github.com/ze674/EZLine/pkg/services"
	"github.com/ze674/EZLine/templates"
	"log"
	"net/http"
	"strconv"
)

// TaskHandler обрабатывает запросы, связанные с заданиями
type TaskHandler struct {
	taskService *services.TaskService
	scanService *services.ScanService
}

// NewTaskHandler создает новый обработчик заданий
func NewTaskHandler(taskService *services.TaskService, scanService *services.ScanService) *TaskHandler {
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

// StartTaskHandler запускает задание на выполнение
func (h *TaskHandler) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID задания", http.StatusBadRequest)
		return
	}

	// Пытаемся запустить задание
	err = h.taskService.StartTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Перенаправляем на страницу сканирования
	http.Redirect(w, r, "/scanning/view", http.StatusSeeOther)
}

// ViewScanningHandler отображает интерфейс сканирования для активного задания
func (h *TaskHandler) ViewScanningHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что есть активное задание
	if !h.taskService.HasActiveTask() {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Получаем информацию о текущем задании
	task := h.taskService.GetActiveTask()
	if task == nil {
		http.Error(w, "Ошибка при получении информации о задании", http.StatusInternalServerError)
		return
	}

	// Получаем статус сканера
	connectionStatus := h.scanService.GetConnectionStatus()

	// Проверяем, есть ли активный ролик
	hasActiveRoll := h.scanService.HasActiveRoll()

	// Получаем статистику сканирования
	var stats *models.ScanStats
	currentRollNumber := 0

	if h.scanService.IsConnected() {
		s := h.scanService.GetStats()
		stats = &s
		currentRollNumber = h.scanService.GetCurrentRollNumber()
	}

	// Рендерим шаблон сканирования
	component := templates.ScanningScreen(*task, connectionStatus, hasActiveRoll, stats, currentRollNumber)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// StartRollHandler обрабатывает запрос на начало сканирования ролика
func (h *TaskHandler) StartRollHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что есть активное задание
	if !h.taskService.HasActiveTask() {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Пытаемся начать сканирование нового ролика
	err := h.scanService.StartNewRoll()
	if err != nil {
		http.Error(w, "Ошибка при создании ролика: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу сканирования
	http.Redirect(w, r, "/scanning/view", http.StatusSeeOther)
}

// FinishRollHandler обрабатывает запрос на завершение сканирования ролика
func (h *TaskHandler) FinishRollHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что есть активное задание
	if !h.taskService.HasActiveTask() {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Завершаем текущий ролик
	err := h.scanService.FinishCurrentRoll()
	if err != nil {
		log.Printf("Ошибка при закрытии ролика: %v", err)
	}

	// Перенаправляем на страницу сканирования
	http.Redirect(w, r, "/scanning/view", http.StatusSeeOther)
}

// RefreshStatsHandler обрабатывает запрос на обновление статистики сканирования
func (h *TaskHandler) RefreshStatsHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что есть активное задание
	if !h.taskService.HasActiveTask() {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Перенаправляем на представление сканирования для обновления
	h.ViewScanningHandler(w, r)
}

// FinishTaskHandler обрабатывает запрос на завершение задания
func (h *TaskHandler) FinishTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Останавливаем сканирование, если оно активно
	if h.scanService.IsScanning() {
		h.scanService.StopScanning()
	}

	// Если есть активный ролик, завершаем его
	if h.scanService.HasActiveRoll() {
		err := h.scanService.FinishCurrentRoll()
		if err != nil {
			log.Printf("Ошибка при закрытии ролика: %v", err)
		}
	}

	// Завершаем задание
	err := h.taskService.FinishTask()
	if err != nil {
		log.Printf("Ошибка при завершении задания: %v", err)
	}

	// Отключаем сканер
	h.scanService.Disconnect()

	// Перенаправляем на список заданий
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
