package validator

import (
	"errors"
	"strings"
)

// Константы ошибок
var (
	ErrInvalidGTIN   = errors.New("GTIN не найден в коде")
	ErrInvalidLength = errors.New("код имеет некорректную длину")
)

// Результат валидации
type ValidationResult struct {
	Valid   bool
	Message string
	Code    string
}

// CodeValidator проверяет штрих-коды на соответствие требованиям
type CodeValidator struct {
	GTIN       string // Код GTIN продукта, который должен содержаться в штрих-коде
	LengthCode int    // Ожидаемая длина штрих-кода
}

// NewCodeValidator создает новый экземпляр валидатора
func NewCodeValidator(gtin string, lengthCode int) *CodeValidator {
	return &CodeValidator{
		GTIN:       gtin,
		LengthCode: lengthCode,
	}
}

// ValidateCode проверяет код и возвращает результат валидации
func (v *CodeValidator) ValidateCode(code string) ValidationResult {
	result := ValidationResult{
		Valid: true,
		Code:  code,
	}

	// Проверка длины кода
	if len(code) != v.LengthCode {
		result.Valid = false
		result.Message = ErrInvalidLength.Error()
		return result
	}

	// Проверка на содержание правильного GTIN
	if !strings.Contains(code, v.GTIN) {
		result.Valid = false
		result.Message = ErrInvalidGTIN.Error()
		return result
	}

	result.Message = "Код валиден"
	return result
}
