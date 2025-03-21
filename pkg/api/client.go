package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Response общая структура ответа от API
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
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
func (c *FactoryClient) GetTasks(lineID int) ([]Task, error) {
	url := fmt.Sprintf("%s/api/tasks?line_id=%d", c.BaseURL, lineID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("ошибка API: %s", response.Error)
	}

	// Преобразуем data в массив задач
	tasksData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршалинге данных: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(tasksData, &tasks); err != nil {
		return nil, fmt.Errorf("ошибка при демаршалинге заданий: %w", err)
	}

	return tasks, nil
}

// Task представляет производственное задание
type Task struct {
	ID          int       `json:"ID"`          // Уникальный идентификатор
	ProductID   int       `json:"ProductID"`   // ID продукта
	ProductName string    `json:"ProductName"` // Название продукта
	LineID      int       `json:"LineID"`      // ID производственной линии
	LineName    string    `json:"LineName"`    // Название линии
	Date        string    `json:"Date"`        // Дата в формате ДД.ММ.ГГГГ
	BatchNumber string    `json:"BatchNumber"` // Номер партии
	Status      string    `json:"Status"`      // Статус задания
	CreatedAt   time.Time `json:"CreatedAt"`   // Дата создания
}
