package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SerialGenerator отвечает за генерацию серийных номеров для коробок
type SerialGenerator struct {
	mu          sync.Mutex
	lastSerial  int             // последний сгенерированный серийный номер
	initialized bool            // флаг инициализации
	gtin        string          // GTIN продукта
	date        string          // дата задания
	batchNumber string          // номер партии
	storagePath string          // путь для хранения файлов
	taskID      int             // ID текущего задания
	codesFile   *os.File        // файл для записи кодов
	usedCodes   map[string]bool // множество уже использованных кодов
}

// NewSerialGenerator создает новый генератор серийных номеров
func NewSerialGenerator(storagePath string) *SerialGenerator {
	// Создаем директорию для хранения, если она не существует
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(storagePath, 0755)
	}

	return &SerialGenerator{
		lastSerial:  0,
		initialized: false,
		storagePath: storagePath,
		usedCodes:   make(map[string]bool),
	}
}

// Initialize инициализирует генератор для конкретного задания
func (g *SerialGenerator) Initialize(taskID int, gtin, date, batchNumber string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.taskID = taskID
	g.gtin = gtin
	g.date = strings.ReplaceAll(date, ".", "") // Удаляем точки из даты
	g.batchNumber = batchNumber
	g.initialized = true
	g.usedCodes = make(map[string]bool) // Инициализируем пустым множеством

	// Создаем файл для записи кодов, если его еще нет
	codesFilename := filepath.Join(g.storagePath, fmt.Sprintf("task_%d_codes.txt", g.taskID))
	var err error

	// Проверяем, существует ли файл с кодами
	if _, err := os.Stat(codesFilename); err == nil {
		// Файл существует, загружаем коды
		err = g.loadExistingCodes(codesFilename)
		if err != nil {
			return fmt.Errorf("ошибка при загрузке существующих кодов: %w", err)
		}
	}

	// Открываем файл для добавления новых кодов
	g.codesFile, err = os.OpenFile(codesFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла кодов: %w", err)
	}

	// Загружаем последний серийный номер из файла
	return g.loadLastSerial()
}

// Метод для загрузки существующих кодов из файла
func (g *SerialGenerator) loadExistingCodes(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла для чтения: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Проверяем, является ли строка серийным номером (не начинается с пробелов)
		if !strings.HasPrefix(line, " ") && len(line) > 0 {
			continue
		} else if strings.HasPrefix(line, "   ") && len(line) > 3 {
			// Эта строка содержит код (начинается с пробелов)
			code := strings.TrimSpace(line)
			if code != "" {
				// Добавляем код в множество использованных
				g.usedCodes[code] = true
				fmt.Printf("Загружен код: %s\n", code)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	fmt.Printf("Загружено %d уникальных кодов из файла\n", len(g.usedCodes))
	return nil
}

// IsCodeUsed проверяет, был ли код уже использован
func (g *SerialGenerator) IsCodeUsed(code string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.usedCodes[code]
}

// MarkCodesAsUsed помечает коды как использованные
func (g *SerialGenerator) MarkCodesAsUsed(codes []string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, code := range codes {
		g.usedCodes[code] = true
	}
}

// GenerateBoxSerial генерирует новый серийный номер для коробки
func (g *SerialGenerator) GenerateBoxSerial() (string, string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.initialized {
		return "", "", fmt.Errorf("генератор не инициализирован")
	}

	// Увеличиваем счетчик серийного номера
	g.lastSerial++

	// Форматируем серийный номер как 6-значное число с ведущими нулями
	serialPart := fmt.Sprintf("%06d", g.lastSerial)

	// Собираем полный серийный номер коробки
	boxSerial := g.gtin + g.date + g.batchNumber + serialPart

	// Сохраняем новое значение в файл
	err := g.saveLastSerial()
	if err != nil {
		return "", "", fmt.Errorf("ошибка при сохранении серийного номера: %w", err)
	}

	return boxSerial, serialPart, nil
}

// SaveCodes сохраняет серийный номер и коды в файл
func (g *SerialGenerator) SaveCodes(boxSerial string, codes []string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.initialized || g.codesFile == nil {
		return fmt.Errorf("генератор не инициализирован или файл кодов не открыт")
	}

	// Добавляем временную метку для удобства отслеживания
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Форматируем запись: серийный номер, временная метка и коды
	record := fmt.Sprintf("%s [%s]\n", boxSerial, timestamp)
	for _, code := range codes {
		record += fmt.Sprintf("   %s\n", code)
		// Отмечаем код как использованный
		g.usedCodes[code] = true
	}
	record += "\n" // Добавляем пустую строку для разделения записей

	// Записываем в файл
	_, err := g.codesFile.WriteString(record)
	if err != nil {
		return fmt.Errorf("ошибка при записи кодов в файл: %w", err)
	}

	// Сбрасываем буфер записи
	err = g.codesFile.Sync()
	if err != nil {
		return fmt.Errorf("ошибка при синхронизации файла: %w", err)
	}

	return nil
}

// Close закрывает все открытые файлы
func (g *SerialGenerator) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.codesFile != nil {
		err := g.codesFile.Close()
		g.codesFile = nil
		if err != nil {
			return fmt.Errorf("ошибка при закрытии файла кодов: %w", err)
		}
	}

	return nil
}

// GetLastSerial возвращает последний сгенерированный серийный номер
func (g *SerialGenerator) GetLastSerial() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.lastSerial
}

// Метод для сохранения последнего серийного номера в файл
func (g *SerialGenerator) saveLastSerial() error {
	if g.taskID == 0 {
		return fmt.Errorf("не указан ID задания")
	}

	filename := filepath.Join(g.storagePath, fmt.Sprintf("task_%d_serial.txt", g.taskID))
	return os.WriteFile(filename, []byte(strconv.Itoa(g.lastSerial)), 0644)
}

// Метод для загрузки последнего серийного номера из файла
func (g *SerialGenerator) loadLastSerial() error {
	if g.taskID == 0 {
		return fmt.Errorf("не указан ID задания")
	}

	filename := filepath.Join(g.storagePath, fmt.Sprintf("task_%d_serial.txt", g.taskID))

	// Проверяем, существует ли файл
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Файл не существует, начинаем с 0
		g.lastSerial = 0
		return nil
	}

	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла серийного номера: %w", err)
	}

	// Преобразуем данные в число
	lastSerial, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("ошибка при преобразовании серийного номера: %w", err)
	}

	g.lastSerial = lastSerial
	return nil
}
