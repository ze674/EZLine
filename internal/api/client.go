package api

import (
	"encoding/json"
	"fmt"
	"github.com/ze674/EZLine/internal/models"
	"net/http"
	"net/url"
	"time"
)

// Response общая структура ответа от API
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error,omitempty"`
}

// FactoryClient предоставляет методы для взаимодействия с API EZFactory
type FactoryClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewFactoryClient создает новый экземпляр клиента для работы с EZFactory
func NewFactoryClient(baseURL string) *FactoryClient {
	return &FactoryClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTasks получает список заданий для указанной линии
func (c *FactoryClient) GetTasks(lineID int) ([]models.Task, error) {
	urlStr := fmt.Sprintf("%s/api/tasks?line_id=%d", c.BaseURL, lineID)

	tasks := make([]models.Task, 0)

	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return tasks, fmt.Errorf("ошибка при запросе к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tasks, fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return tasks, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if !response.Success {
		return tasks, fmt.Errorf("ошибка API: %s", response.Error)
	}

	if err := json.Unmarshal(response.Data, &tasks); err != nil {
		return []models.Task{}, fmt.Errorf("ошибка при демаршалинге заданий: %w", err)
	}

	return tasks, nil
}

// GetTaskByID получает информацию о задании по ID
func (c *FactoryClient) GetTaskByID(taskID int) (models.Task, error) {
	urlStr := fmt.Sprintf("%s/api/tasks/%d", c.BaseURL, taskID)

	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return models.Task{}, fmt.Errorf("ошибка при запросе к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Task{}, fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return models.Task{}, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if !response.Success {
		return models.Task{}, fmt.Errorf("ошибка API: %s", response.Error)
	}

	var task models.Task
	if err := json.Unmarshal(response.Data, &task); err != nil {
		return models.Task{}, fmt.Errorf("ошибка при демаршалинге задания: %w", err)
	}

	return task, nil
}

// UpdateTaskStatus обновляет статус задания
func (c *FactoryClient) UpdateTaskStatus(taskID int, newStatus string) error {
	urlStr := fmt.Sprintf("%s/api/tasks/%d/status", c.BaseURL, taskID)

	// Подготовка данных формы
	data := url.Values{}
	data.Set("status", newStatus)

	// Отправка POST запроса
	resp, err := c.HTTPClient.PostForm(urlStr, data)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	return nil
}
