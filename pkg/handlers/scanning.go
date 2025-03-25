package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/ze674/EZLine/pkg/models"
	"github.com/ze674/EZLine/pkg/scanner"
	"github.com/ze674/EZLine/templates"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	tcpScanner     *scanner.TCPScanner
	rollManager    *scanner.RollManager
	scannerAddress string
	storagePath    string
	codeLength     int

	currentTask models.Task

	//Мьютекс для синхронизации доступа к глобальным переменным
	scannerMutex sync.Mutex
)

// StartScanningHandler обрабатывает запрос на начало сканирования
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

	// Получаем информацию о задании через API
	task, err := FactoryClient.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(), http.StatusInternalServerError)
		return
	}

	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Устанавливаем задание как активное
	ActiveTaskID = taskID
	currentTask = task

	// Обновляем статус задания на "в работе", если оно еще не в работе
	if task.Status != "в работе" {
		err = FactoryClient.UpdateTaskStatus(taskID, "в работе")
		if err != nil {
			log.Printf("Ошибка при обновлении статуса задания: %v", err)
		}
	}

	// Создаем менеджер роликов
	rollManager, err = scanner.NewRollManager(taskID, storagePath)
	if err != nil {
		http.Error(w, "Ошибка создания менеджера роликов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем подключение к сканеру
	connectionStatus := "Сканер не подключен"

	// Если сканер еще не создан, создаем его
	if tcpScanner == nil {
		tcpScanner = scanner.NewTCPScanner(scannerAddress, codeLength) // Минимальная длина 10, максимальная 50
	}

	// Устанавливаем менеджер роликов для сканера
	tcpScanner.SetRollManager(rollManager)

	// Пробуем подключиться к сканеру, если он еще не подключен
	if !tcpScanner.IsConnected() {
		err = tcpScanner.Connect()
		if err != nil {
			connectionStatus = "Ошибка подключения к сканеру: " + err.Error()
		} else {
			connectionStatus = "Сканер успешно подключен"
		}
	} else {
		connectionStatus = "Сканер уже подключен"
	}

	// Проверяем, есть ли активный ролик
	hasActiveRoll := rollManager.HasActiveRoll()

	// Отображаем экран сканирования
	component := templates.ScanningScreen(task, connectionStatus, hasActiveRoll, nil, 0)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// StartTaskHandler - запускает новое задание
func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID задания из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID задания", http.StatusBadRequest)
		return
	}

	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Если уже есть активное задание, не даем начать новое
	if ActiveTaskID != 0 && ActiveTaskID != taskID {
		http.Error(w, "Уже выполняется другое задание. Завершите его перед началом нового.", http.StatusBadRequest)
		return
	}

	// Если это уже активное задание, просто перенаправляем на сканирование
	if ActiveTaskID == taskID && rollManager != nil && rollManager.GetTaskID() == taskID {
		http.Redirect(w, r, "/scanning/view", http.StatusSeeOther)
		return
	}

	// Получаем информацию о задании через API
	task, err := FactoryClient.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Ошибка при получении информации о задании: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем задание как активное
	ActiveTaskID = taskID
	currentTask = task

	// Обновляем статус задания на "в работе", если оно еще не в работе
	if task.Status != "в работе" {
		err = FactoryClient.UpdateTaskStatus(taskID, "в работе")
		if err != nil {
			log.Printf("Ошибка при обновлении статуса задания: %v", err)
		}
	}

	// Создаем менеджер роликов
	rollManager, err = scanner.NewRollManager(taskID, storagePath)
	if err != nil {
		http.Error(w, "Ошибка создания менеджера роликов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если сканер еще не создан, создаем его
	if tcpScanner == nil {
		tcpScanner = scanner.NewTCPScanner(scannerAddress, codeLength)
	}

	// Устанавливаем менеджер роликов для сканера
	tcpScanner.SetRollManager(rollManager)

	// Перенаправляем на экран сканирования
	http.Redirect(w, r, "/scanning/view", http.StatusSeeOther)
}

// ViewScanningHandler - отображает текущий экран сканирования
func ViewScanningHandler(w http.ResponseWriter, r *http.Request) {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Проверяем, что есть активное задание
	if ActiveTaskID == 0 || &currentTask == nil || rollManager == nil {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Проверяем подключение к сканеру
	connectionStatus := "Сканер не подключен"

	if tcpScanner != nil {
		if !tcpScanner.IsConnected() {
			err := tcpScanner.Connect()
			if err != nil {
				connectionStatus = "Ошибка подключения к сканеру: " + err.Error()
			} else {
				connectionStatus = "Сканер успешно подключен"
			}
		} else {
			connectionStatus = "Сканер уже подключен"
		}
	}

	// Проверяем, есть ли активный ролик
	hasActiveRoll := rollManager.HasActiveRoll()

	// Отображаем экран сканирования
	component := templates.ScanningScreen(currentTask, connectionStatus, hasActiveRoll, nil, 0)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// StartRollHandler обрабатывает запрос на начало сканирования ролика
func StartRollHandler(w http.ResponseWriter, r *http.Request) {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Проверяем, есть ли активное задание
	if ActiveTaskID == 0 {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Проверяем подключение к сканеру
	if tcpScanner == nil || !tcpScanner.IsConnected() {
		http.Error(w, "Сканер не подключен", http.StatusBadRequest)
		return
	}

	// Проверяем, нет ли уже активного ролика
	if rollManager.HasActiveRoll() {
		http.Error(w, "Уже есть активный ролик", http.StatusBadRequest)
		return
	}

	// Начинаем новый ролик
	currentRoll, err := rollManager.StartNewRoll()
	if err != nil {
		http.Error(w, "Ошибка при создании ролика: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Начинаем прослушивание сканера
	err = tcpScanner.StartListening()
	if err != nil {
		// Если не удалось начать прослушивание, закрываем ролик
		rollManager.FinishCurrentRoll()
		http.Error(w, "Ошибка начала прослушивания сканера: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображаем экран активного сканирования
	stats := tcpScanner.GetStats()
	component := templates.ScanningScreen(currentTask, "Сканер активен", true, &stats, currentRoll)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// FinishRollHandler обрабатывает запрос на завершение сканирования ролика
func FinishRollHandler(w http.ResponseWriter, r *http.Request) {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Проверяем, есть ли активное задание
	if ActiveTaskID == 0 {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Проверяем, есть ли активный ролик
	if rollManager == nil || !rollManager.HasActiveRoll() {
		http.Error(w, "Нет активного ролика", http.StatusBadRequest)
		return
	}

	// Останавливаем прослушивание сканера
	if tcpScanner != nil {
		tcpScanner.StopListening()
	}

	// Завершаем текущий ролик
	err := rollManager.FinishCurrentRoll()
	if err != nil {
		log.Printf("Ошибка при закрытии ролика: %v", err)
	}

	currentRollNumber := rollManager.GetCurrentRoll()

	// Отображаем экран сканирования
	stats := tcpScanner.GetStats()
	component := templates.ScanningScreen(currentTask, "Сканер подключен", false, &stats, currentRollNumber)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// RefreshStatsHandler обрабатывает запрос на обновление статистики сканирования
func RefreshStatsHandler(w http.ResponseWriter, r *http.Request) {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Проверяем, есть ли активное задание
	if ActiveTaskID == 0 {
		http.Error(w, "Нет активного задания", http.StatusBadRequest)
		return
	}

	// Получаем статистику сканирования
	var stats *models.ScanStats
	if tcpScanner != nil {
		s := tcpScanner.GetStats()
		stats = &s
	}

	// Определяем статус сканера
	connectionStatus := "Сканер не подключен"
	if tcpScanner != nil && tcpScanner.IsConnected() {
		connectionStatus = "Сканер подключен"
	}

	// Определяем наличие активного ролика
	hasActiveRoll := rollManager != nil && rollManager.HasActiveRoll()
	// Получаем номер текущего ролика
	currentRollNumber := 0
	if rollManager != nil {
		currentRollNumber = rollManager.GetCurrentRoll()
	}

	// Отображаем экран сканирования с обновленной статистикой
	component := templates.ScanningScreen(currentTask, connectionStatus, hasActiveRoll, stats, currentRollNumber)

	if r.Header.Get("HX-Request") == "true" {
		component.Render(r.Context(), w)
	} else {
		templates.Page(component).Render(r.Context(), w)
	}
}

// FinishScanningHandler обрабатывает запрос на завершение сканирования
func FinishScanningHandler(w http.ResponseWriter, r *http.Request) {
	scannerMutex.Lock()
	defer scannerMutex.Unlock()

	// Если есть активный ролик, завершаем его
	if rollManager != nil && rollManager.HasActiveRoll() {
		// Останавливаем прослушивание сканера
		if tcpScanner != nil {
			tcpScanner.StopListening()
		}

		// Завершаем ролик
		err := rollManager.FinishCurrentRoll()
		if err != nil {
			log.Printf("Ошибка при закрытии ролика: %v", err)
		}
	}

	// Если сканер подключен, отключаем его
	if tcpScanner != nil && tcpScanner.IsConnected() {
		err := tcpScanner.Disconnect()
		if err != nil {
			log.Printf("Ошибка при отключении сканера: %v", err)
		}
	}

	// Обновляем статус задания на "завершено"
	if ActiveTaskID != 0 {
		err := FactoryClient.UpdateTaskStatus(ActiveTaskID, "завершено")
		if err != nil {
			log.Printf("Ошибка при обновлении статуса задания: %v", err)
		}
	}

	// Сбрасываем глобальные переменные
	tcpScanner = nil
	rollManager = nil
	ActiveTaskID = 0

	// Перенаправляем на список заданий
	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
