package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Проверка аргументов командной строки
	if len(os.Args) < 5 {
		fmt.Println("Использование: program input.txt output.txt initial_serial_number box_size")
		fmt.Println("  input.txt - входной файл с кодами (по одному на строку)")
		fmt.Println("  output.txt - выходной XML-файл")
		fmt.Println("  initial_serial_number - начальный серийный номер короба (например, 01146070547613261125021210001702100001)")
		fmt.Println("  box_size - количество кодов в одном коробе")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	initialSerialNumber := os.Args[3]
	boxSize, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Ошибка при преобразовании размера короба: %v\n", err)
		return
	}

	// Проверка размера короба
	if boxSize <= 0 {
		fmt.Println("Размер короба должен быть положительным числом")
		return
	}

	// Открываем входной файл
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Ошибка при открытии входного файла: %v\n", err)
		return
	}
	defer input.Close()

	// Создаем выходной файл
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Ошибка при создании выходного файла: %v\n", err)
		return
	}
	defer output.Close()

	// Начало XML-файла
	output.WriteString("<root> \n")

	scanner := bufio.NewScanner(input)
	currentSerialNumber := initialSerialNumber
	codesInCurrentBox := 0

	// Читаем коды из входного файла
	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки
		if code == "" {
			continue
		}

		// Проверяем, нужно ли увеличить серийный номер короба
		if codesInCurrentBox >= boxSize {
			// Увеличиваем серийный номер короба на 1
			currentSerialNumber = incrementSerialNumber(currentSerialNumber)
			codesInCurrentBox = 0
		}

		// Формируем строку XML
		xmlLine := fmt.Sprintf("    <CodeNamesSerial><cCodeNameSerial>%s</cCodeNameSerial><cCodeNameSerialParent>%s</cCodeNameSerialParent><cOutID>WMS104388</cOutID></CodeNamesSerial>\n",
			code, currentSerialNumber)

		// Записываем в выходной файл
		output.WriteString(xmlLine)

		// Увеличиваем счетчик кодов в текущем коробе
		codesInCurrentBox++
	}

	// Проверяем ошибки сканирования
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении входного файла: %v\n", err)
	}

	// Конец XML-файла
	output.WriteString("</root>")

	fmt.Println("Преобразование успешно завершено!")
}

// Функция для увеличения серийного номера короба на 1
func incrementSerialNumber(serialNumber string) string {
	// Преобразуем строку в число (если это возможно)
	// Если невозможно, просто добавляем "1" в конец строки как костыль

	// Предполагаем, что последние несколько цифр - это инкрементируемая часть
	// В этом примере берем последние 5 цифр, но можно настроить по необходимости

	length := len(serialNumber)
	if length < 5 {
		// Слишком короткий серийный номер, просто добавляем 1
		return serialNumber + "1"
	}

	// Получаем последние 5 цифр
	lastPart := serialNumber[length-5:]
	prefixPart := serialNumber[:length-5]

	// Преобразуем последнюю часть в число
	lastPartNum, err := strconv.Atoi(lastPart)
	if err != nil {
		// Если не удается преобразовать, просто увеличиваем последний символ
		return serialNumber[:length-1] + string(serialNumber[length-1]+1)
	}

	// Увеличиваем число на 1
	lastPartNum++

	// Форматируем обратно в строку с ведущими нулями
	newLastPart := fmt.Sprintf("%05d", lastPartNum)

	// Собираем новый серийный номер
	return prefixPart + newLastPart
}
