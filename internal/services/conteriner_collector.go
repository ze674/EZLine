// internal/services/container_collector.go
package services

import (
	"sync"
)

// ContainerCollector отвечает за накопление кодов продуктов до заполнения контейнера
type ContainerCollector struct {
	mu                sync.Mutex
	pendingCodes      []string // Накопленные коды
	containerCapacity int      // Емкость короба
	layerCapacity     int      // Емкость слоя
	totalLayers       int      // Общее количество слоев
}

// NewContainerCollector создает новый сервис сбора контейнера
func NewContainerCollector(containerCapacity, layerCapacity, totalLayers int) *ContainerCollector {
	return &ContainerCollector{
		pendingCodes:      []string{},
		containerCapacity: containerCapacity,
		layerCapacity:     layerCapacity,
		totalLayers:       totalLayers,
	}
}

// AddCodes добавляет новые коды в накопление
// Возвращает: есть ли дубликаты, текущее количество, достигли ли емкости короба
func (c *ContainerCollector) AddCodes(newCodes []string) (bool, int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверка на превышение емкости слоя
	if len(newCodes) > c.layerCapacity {
		return false, len(c.pendingCodes), false
	}

	// Проверка на дубликаты
	hasDuplicates := false
	for _, newCode := range newCodes {
		for _, pendingCode := range c.pendingCodes {
			if newCode == pendingCode {
				hasDuplicates = true
				return hasDuplicates, len(c.pendingCodes), false
			}
		}
	}

	// Добавляем коды
	c.pendingCodes = append(c.pendingCodes, newCodes...)

	// Проверяем, достигли ли емкости короба
	isBoxFull := len(c.pendingCodes) >= c.containerCapacity

	return hasDuplicates, len(c.pendingCodes), isBoxFull
}

// GetPendingCodes возвращает текущие накопленные коды
func (c *ContainerCollector) GetPendingCodes() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем копию, чтобы избежать проблем с параллельным доступом
	result := make([]string, len(c.pendingCodes))
	copy(result, c.pendingCodes)

	return result
}

// ExtractCodes извлекает накопленные коды для сохранения
// Возвращает коды и сбрасывает внутреннее состояние
func (c *ContainerCollector) ExtractCodes() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.pendingCodes) == 0 {
		return nil
	}

	// Ограничиваем количество возвращаемых кодов емкостью короба
	extractCount := len(c.pendingCodes)
	if extractCount > c.containerCapacity {
		extractCount = c.containerCapacity
	}

	// Извлекаем коды
	result := make([]string, extractCount)
	copy(result, c.pendingCodes[:extractCount])

	// Сбрасываем накопленные коды
	c.pendingCodes = []string{}

	return result
}

// Reset сбрасывает все накопленные коды
func (c *ContainerCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingCodes = []string{}
}

// GetBoxCapacity возвращает емкость короба
func (c *ContainerCollector) GetContainerCapacity() int {
	return c.containerCapacity
}

// GetLayerCapacity возвращает емкость слоя
func (c *ContainerCollector) GetLayerCapacity() int {
	return c.layerCapacity
}

// IsFull проверяет, достигли ли мы емкости короба
func (c *ContainerCollector) IsFull() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.pendingCodes) >= c.containerCapacity
}

// GetCount возвращает текущее количество накопленных кодов
func (c *ContainerCollector) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.pendingCodes)
}
