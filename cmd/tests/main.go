package main

import (
	"fmt"
	"strconv"
)

// GenerateITF14 конвертирует EAN-13 в ITF-14
func GenerateITF14(ean13 string) (string, error) {
	// Проверяем длину EAN-13
	if len(ean13) != 13 {
		return "", fmt.Errorf("EAN-13 должен содержать 13 цифр")
	}

	// Убираем последнюю контрольную цифру из EAN-13
	baseCode := ean13[:12]

	// Добавляем начальную цифру упаковки (1)
	itf14 := "1" + baseCode

	// Вычисляем контрольную цифру
	checksum := calculateITF14Checksum(itf14)

	// Добавляем контрольную цифру
	itf14 += strconv.Itoa(checksum)

	return itf14, nil
}

// calculateITF14Checksum вычисляет контрольную цифру для ITF-14
func calculateITF14Checksum(code string) int {
	// Статические веса для ITF-14
	weights := []int{3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1}

	// Сумма произведений цифр на веса
	total := 0
	for i := 0; i < len(code); i++ {
		digit, _ := strconv.Atoi(string(code[i]))
		total += digit * weights[i]
	}

	// Вычисление контрольной цифры
	checksum := (10 - (total % 10)) % 10
	return checksum
}

func main() {
	// Пример использования
	ean13 := "4601234567890"
	itf14, err := GenerateITF14(ean13)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("EAN-13:", ean13)
	fmt.Println("ITF-14:", itf14)

	// Дополнительные тесты
	testCases := []string{
		"4601234567890",
		"1234567890123",
		"9876543210987",
		"4607054761244",
	}

	for _, testEAN := range testCases {
		result, err := GenerateITF14(testEAN)
		if err != nil {
			fmt.Printf("Ошибка для %s: %v\n", testEAN, err)
		} else {
			fmt.Printf("EAN-13: %s, ITF-14: %s\n", testEAN, result)
		}
	}
}
