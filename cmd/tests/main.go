package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Использование: program input.txt output.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

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

	// Записываем начало XML
	output.WriteString("<root> \n")

	scanner := bufio.NewScanner(input)
	var currentParent string

	// Построчное чтение входного файла
	for scanner.Scan() {
		line := scanner.Text()

		// Проверяем, является ли строка идентификатором группы
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			// Это строка с идентификатором группы
			parts := strings.Split(line, " [")
			if len(parts) > 0 {
				currentParent = parts[0]
			}
		} else {
			// Это строка с кодом
			code := strings.TrimSpace(line)
			if code != "" {
				// Форматируем и записываем XML-элемент
				xmlElement := fmt.Sprintf("    <CodeNamesSerial><cCodeNameSerial>%s</cCodeNameSerial><cCodeNameSerialParent>%s</cCodeNameSerialParent><cOutID>WMS104388</cOutID></CodeNamesSerial>\n",
					code, currentParent)
				output.WriteString(xmlElement)
			}
		}
	}

	// Записываем конец XML
	output.WriteString("</root>")

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении входного файла: %v\n", err)
	}

	fmt.Println("Преобразование успешно завершено!")
}
