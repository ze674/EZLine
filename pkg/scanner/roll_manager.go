package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var fileSuffix = ".csv"

type RollManager struct {
	taskID      int      //ID задания
	storagePath string   //Базовый путь к хранилищу
	currentRoll int      //Текущий № ролика
	rollFile    *os.File //Текущий открытый файл
}

// NewRollManager создает новый менеджер роликов
func NewRollManager(taskID int, storagePath string) (*RollManager, error) {
	//Создаем менеджер
	rm := &RollManager{
		taskID:      taskID,
		storagePath: storagePath,
		currentRoll: 0,
	}

	//Создаем директорию для задания, если она не существует
	if err := rm.creatreTaskDirectory(); err != nil {
		return nil, err
	}

	return rm, nil
}

// creatreTaskDirectory создает директорию для задания
func (rm *RollManager) creatreTaskDirectory() error {
	taskDir := filepath.Join(rm.storagePath, fmt.Sprintf("task_%d", rm.taskID))
	return os.MkdirAll(taskDir, 0755)
}

// getTaskDirectory возвращает путь к директории задания
func (rm *RollManager) getTaskDirectory() string {
	return filepath.Join(rm.storagePath, fmt.Sprintf("task_%d", rm.taskID))
}

// getNextRollNumber возвращает следующий № ролика
func (rm *RollManager) getNextRollNumber() (int, error) {
	taskDir := rm.getTaskDirectory()

	//Проверяем существуют ли файлы
	files, err := os.ReadDir(taskDir)
	if err != nil {
		return 0, err
	}

	maxRoll := 0

	//Ищем файлы с именем вида "1.txt", "2.txt" и т.д.
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !strings.HasSuffix(name, fileSuffix) {
			continue
		}

		//Получаем номер ролика из имени файла
		numStr := strings.TrimSuffix(name, fileSuffix)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue //Пропускаем если имя не является числом
		}

		if num > maxRoll {
			maxRoll = num
		}
	}

	// Следующий номер ролика = максимальный найденный + 1
	return maxRoll + 1, nil
}

// StartNewRoll начинает новый ролик
func (rm *RollManager) StartNewRoll() (int, error) {
	//Закрываем предыдущий файл, если он был открыт
	if rm.rollFile != nil {
		rm.rollFile.Close()
		rm.rollFile = nil
	}

	//Получаем следующий номер ролика
	nextRoll, err := rm.getNextRollNumber()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении следующего ролика: %w", err)
	}

	//Создаем файл для нового ролика
	rollPath := filepath.Join(rm.getTaskDirectory(), fmt.Sprintf("%d%s", nextRoll, fileSuffix))
	file, err := os.Create(rollPath)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании файла ролика: %w", err)
	}

	rm.rollFile = file
	rm.currentRoll = nextRoll

	return nextRoll, nil
}

// FinishCurrentRoll завершает текущий ролик
func (rm *RollManager) FinishCurrentRoll() error {
	if rm.rollFile == nil {
		return nil //Ничего не делаем, если нет активного ролика
	}

	err := rm.rollFile.Close()
	rm.rollFile = nil
	return err
}

// WriteCode записывает код в текущий ролик
func (rm *RollManager) WriteCode(code string) error {
	if rm.rollFile == nil {
		return fmt.Errorf("нет активного ролика")
	}

	_, err := rm.rollFile.WriteString(code + "\n")
	return err
}

// GetCurrentRoll возвращает № текущего ролика
func (rm *RollManager) GetCurrentRoll() int {
	return rm.currentRoll
}

// HasActiveRoll возвращает true, если есть активный ролик
func (rm *RollManager) HasActiveRoll() bool {
	return rm.rollFile != nil
}

func (rm *RollManager) GetTaskID() int {
	return rm.taskID
}
