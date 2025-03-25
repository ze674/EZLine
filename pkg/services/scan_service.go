package services

import (
	"fmt"
	"github.com/ze674/EZLine/pkg/models"
	"github.com/ze674/EZLine/pkg/scanner"
	"log"
	"sync"
)

// ScanService представляет сервис для работы со сканером и роликами
type ScanService struct {
	mu               sync.Mutex
	scanner          *scanner.TCPScanner
	rollManager      *scanner.RollManager
	scannerAddress   string
	storagePath      string
	codeLength       int
	isConnected      bool
	isScanning       bool
	connectionStatus string
}

// NewScanService создает новый сервис для работы со сканером
func NewScanService(scannerAddress, storagePath string, codeLength int) *ScanService {
	return &ScanService{
		scannerAddress:   scannerAddress,
		storagePath:      storagePath,
		codeLength:       codeLength,
		isConnected:      false,
		isScanning:       false,
		connectionStatus: "Сканер не подключен",
	}
}

// Initialize инициализирует сканер для задания
func (s *ScanService) Initialize(taskID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Создаем менеджер роликов
	rollManager, err := scanner.NewRollManager(taskID, s.storagePath)
	if err != nil {
		return fmt.Errorf("ошибка создания менеджера роликов: %w", err)
	}
	s.rollManager = rollManager

	// Создаем сканер, если он еще не создан
	if s.scanner == nil {
		s.scanner = scanner.NewTCPScanner(s.scannerAddress, s.codeLength)
	}

	// Устанавливаем менеджер роликов для сканера
	s.scanner.SetRollManager(rollManager)

	return nil
}

// Connect подключается к сканеру
func (s *ScanService) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanner == nil {
		s.scanner = scanner.NewTCPScanner(s.scannerAddress, s.codeLength)
	}

	if !s.scanner.IsConnected() {
		err := s.scanner.Connect()
		if err != nil {
			s.connectionStatus = "Ошибка подключения к сканеру: " + err.Error()
			s.isConnected = false
			return err
		}
		s.connectionStatus = "Сканер успешно подключен"
		s.isConnected = true
	} else {
		s.connectionStatus = "Сканер уже подключен"
		s.isConnected = true
	}

	return nil
}

// Disconnect отключается от сканера
func (s *ScanService) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanner != nil && s.scanner.IsConnected() {
		// Сначала останавливаем сканирование, если оно активно
		if s.isScanning {
			s.scanner.StopListening()
			s.isScanning = false
		}

		err := s.scanner.Disconnect()
		if err != nil {
			log.Printf("Ошибка при отключении сканера: %v", err)
			return err
		}
	}

	s.isConnected = false
	s.connectionStatus = "Сканер отключен"
	return nil
}

// IsConnected проверяет, подключен ли сканер
func (s *ScanService) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isConnected
}

// IsScanning проверяет, идет ли процесс сканирования
func (s *ScanService) IsScanning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isScanning
}

// GetConnectionStatus возвращает статус подключения сканера
func (s *ScanService) GetConnectionStatus() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.connectionStatus
}

// StartNewRoll начинает новый ролик и запускает сканирование
func (s *ScanService) StartNewRoll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, подключен ли сканер
	if s.scanner == nil || !s.scanner.IsConnected() {
		return fmt.Errorf("сканер не подключен")
	}

	// Проверяем, есть ли менеджер роликов
	if s.rollManager == nil {
		return fmt.Errorf("менеджер роликов не инициализирован")
	}

	// Проверяем, нет ли уже активного ролика
	if s.rollManager.HasActiveRoll() {
		return fmt.Errorf("уже есть активный ролик")
	}

	// Начинаем новый ролик
	_, err := s.rollManager.StartNewRoll()
	if err != nil {
		return fmt.Errorf("ошибка при создании ролика: %w", err)
	}

	// Начинаем прослушивание сканера
	err = s.scanner.StartListening()
	if err != nil {
		// Если не удалось начать прослушивание, закрываем ролик
		s.rollManager.FinishCurrentRoll()
		return fmt.Errorf("ошибка при начале прослушивания сканера: %w", err)
	}

	s.isScanning = true
	s.connectionStatus = "Сканер активен"
	return nil
}

// FinishCurrentRoll завершает текущий ролик и останавливает сканирование
func (s *ScanService) FinishCurrentRoll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, есть ли активный ролик
	if s.rollManager == nil || !s.rollManager.HasActiveRoll() {
		return fmt.Errorf("нет активного ролика")
	}

	// Останавливаем прослушивание сканера
	if s.scanner != nil {
		s.scanner.StopListening()
		s.isScanning = false
	}

	// Завершаем текущий ролик
	err := s.rollManager.FinishCurrentRoll()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии ролика: %w", err)
	}

	s.connectionStatus = "Сканер подключен"
	return nil
}

// StopScanning останавливает процесс сканирования
func (s *ScanService) StopScanning() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanner != nil {
		s.scanner.StopListening()
	}
	s.isScanning = false
}

// HasActiveRoll проверяет, есть ли активный ролик
func (s *ScanService) HasActiveRoll() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rollManager != nil && s.rollManager.HasActiveRoll()
}

// GetCurrentRollNumber возвращает номер текущего ролика
func (s *ScanService) GetCurrentRollNumber() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.rollManager == nil {
		return 0
	}
	return s.rollManager.GetCurrentRoll()
}

// GetStats возвращает статистику сканирования
func (s *ScanService) GetStats() models.ScanStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.scanner == nil {
		return models.ScanStats{}
	}
	return s.scanner.GetStats()
}
