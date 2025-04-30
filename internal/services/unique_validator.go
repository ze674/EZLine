// internal/services/code_validator.go
package services

import (
	"fmt"
	"github.com/ze674/EZLine/internal/repository"
	"sync"
)

// CodeUniquenessValidator отвечает за проверку уникальности кодов в пределах задания
type CodeUniquenessValidator struct {
	mu             sync.Mutex
	taskID         int
	initialized    bool
	usedCodes      map[string]bool // Кэш использованных кодов
	itemRepository *repository.ItemRepository
}

// NewCodeUniquenessValidator создает новый валидатор уникальности кодов
func NewCodeUniquenessValidator() *CodeUniquenessValidator {
	return &CodeUniquenessValidator{
		itemRepository: repository.NewItemRepository(),
		usedCodes:      make(map[string]bool),
		initialized:    false,
	}
}

// Initialize инициализирует валидатор для конкретного задания
func (v *CodeUniquenessValidator) Initialize(taskID int) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.taskID = taskID
	v.usedCodes = make(map[string]bool)

	// Загружаем все существующие коды для задания
	items, err := v.itemRepository.GetItemsByTaskID(taskID)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке списка товаров: %w", err)
	}

	// Заполняем кэш
	for _, item := range items {
		v.usedCodes[item.Code] = true
	}

	v.initialized = true
	return nil
}

// IsCodeUnique проверяет, уникален ли код в пределах задания
func (v *CodeUniquenessValidator) IsCodeUnique(code string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return false // Если не инициализирован, считаем код неуникальным для безопасности
	}

	// Проверяем наличие кода в кэше
	_, exists := v.usedCodes[code]
	return !exists // Если кода нет в кэше, значит он уникальный
}

// IsCodeUsed проверяет, использовался ли уже код
func (v *CodeUniquenessValidator) IsCodeUsed(code string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return false
	}

	_, exists := v.usedCodes[code]
	return exists
}

// IsCodesUnique проверяет, что все коды в списке уникальны в пределах задания
// Возвращает два значения: общий результат проверки и список неуникальных кодов
func (v *CodeUniquenessValidator) IsCodesUnique(codes []string) (bool, []string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return false, codes // Если не инициализирован, считаем все коды неуникальными
	}

	var duplicates []string
	for _, code := range codes {
		if _, exists := v.usedCodes[code]; exists {
			// Код уже использовался
			duplicates = append(duplicates, code)
		}
	}

	// Если список дубликатов пуст, значит все коды уникальны
	return len(duplicates) == 0, duplicates
}

// MarkCodeAsUsed отмечает код как использованный
func (v *CodeUniquenessValidator) MarkCodeAsUsed(code string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return
	}

	v.usedCodes[code] = true
}

// MarkCodesAsUsed отмечает несколько кодов как использованные
func (v *CodeUniquenessValidator) MarkCodesAsUsed(codes []string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !v.initialized {
		return
	}

	for _, code := range codes {
		v.usedCodes[code] = true
	}
}

// GetUsedCodesCount возвращает количество использованных кодов
func (v *CodeUniquenessValidator) GetUsedCodesCount() int {
	v.mu.Lock()
	defer v.mu.Unlock()

	return len(v.usedCodes)
}
