package scanner

import (
	"bufio"
	"fmt"
	"github.com/ze674/EZLine/pkg/models"
	"net"
	"strings"
	"sync"
	"time"
)

var connectTimeout = 5 * time.Second
var readTimeout = 500 * time.Millisecond

var errAlreadyConnected = fmt.Errorf("already connected")

type TCPScanner struct {
	address     string           // Адрес в формате host:port
	conn        net.Conn         // TCP-Соединение
	active      bool             // Флаг активного соединения
	listening   bool             // Флаг активного прослушивания
	mu          sync.Mutex       // Мьютекс для синхронизации
	codes       map[string]bool  //Карта для отслеживания уникальности кодов
	codeLength  int              //Длина кода для валидации
	stats       models.ScanStats //Статистика сканирования
	rollManager *RollManager     //Менеджер роликов
}

// NewTCPScanner создает новый экземпляр TCP-Сканнера
func NewTCPScanner(address string, codeLength int) *TCPScanner {
	return &TCPScanner{
		address:    address,
		active:     false,
		listening:  false,
		codes:      make(map[string]bool),
		codeLength: codeLength,
		stats: models.ScanStats{
			RecentCodes:   make([]string, 0, 10),
			RecentResults: make([]bool, 0, 10),
		},
	}
}

// SetRollManager устанавливает менеджер роликов
func (s *TCPScanner) SetRollManager(rm *RollManager) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rollManager = rm
}

// Connect устанавливает соединение со сканером
func (s *TCPScanner) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если сканер уже подключен, ничего возвращаем ошибку
	if s.active {
		return errAlreadyConnected
	}

	//Устанавливаем соединение
	conn, err := net.DialTimeout("tcp", s.address, connectTimeout)
	if err != nil {
		return fmt.Errorf("can't connect to scanner: %w", err)
	}

	s.conn = conn
	s.active = true
	return nil
}

// Disconnect закрывает соединение со сканером
func (s *TCPScanner) Disconnect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//Останавливаем прослушивание, если оно активно
	s.listening = false

	// Если соединение уже закрыто, ничего возвращаем ошибку
	if !s.active {
		return nil
	}

	s.active = false
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		return err
	}
	return nil
}

// IsConnected возвращает true, если сканер подключен
func (s *TCPScanner) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.active
}

// StartListening начинает прослушивание сканера
func (s *TCPScanner) StartListening() error {
	s.mu.Lock()

	// Если уже прослушиваем или не подключены, возвращаем ошибку
	if s.listening {
		s.mu.Unlock()
		return fmt.Errorf("прослушивание уже запущено")
	}

	if !s.active || s.conn == nil {
		s.mu.Unlock()
		return fmt.Errorf("сканер не подключен")
	}

	// Сбрасываем статистику и карту кодов при начале нового сканирования
	s.stats = models.ScanStats{
		RecentCodes:   make([]string, 0, 10),
		RecentResults: make([]bool, 0, 10),
	}
	s.codes = make(map[string]bool)

	s.listening = true
	s.mu.Unlock()

	// Запускаем горутину для прослушивания
	go s.listen()

	return nil
}

// StopListening останавливает прослушивание сканера
func (s *TCPScanner) StopListening() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listening = false
}

// GetStats возвращает текущую статистику сканирования
func (s *TCPScanner) GetStats() models.ScanStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Создаем копию статистики
	stats := models.ScanStats{
		TotalCodes:    s.stats.TotalCodes,
		ValidCodes:    s.stats.ValidCodes,
		InvalidCodes:  s.stats.InvalidCodes,
		LastCode:      s.stats.LastCode,
		LastResult:    s.stats.LastResult,
		LastError:     s.stats.LastError,
		RecentCodes:   make([]string, len(s.stats.RecentCodes)),
		RecentResults: make([]bool, len(s.stats.RecentResults)),
	}

	copy(stats.RecentCodes, s.stats.RecentCodes)
	copy(stats.RecentResults, s.stats.RecentResults)

	return stats
}

// listen обрабатывает входящие данные из TCP-соединения
func (s *TCPScanner) listen() {
	reader := bufio.NewReader(s.conn)

	for {
		s.mu.Lock()
		if !s.listening || !s.active {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()

		// Устанавливаем таймаут чтения
		s.conn.SetReadDeadline(time.Now().Add(readTimeout))

		// Считываем строку до символа новой строки
		line, err := reader.ReadString('\n')
		if err != nil {
			// Если таймаут, продолжаем
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}

			s.mu.Lock()
			if s.listening && s.active {
				// Добавляем ошибку в статистику
				s.stats.LastCode = ""
				s.stats.LastResult = false
				s.stats.LastError = "Ошибка чтения из TCP: " + err.Error()
			}
			s.listening = false
			s.mu.Unlock()
			return
		}

		// Обрабатываем полученный код
		code := strings.TrimSpace(line)
		s.mu.Lock()

		// Увеличиваем общее количество кодов
		s.stats.TotalCodes++

		// Сохраняем последний код
		s.stats.LastCode = code

		// Проверяем длину кода
		valid := true
		var errMsg string

		if code == "NOREAD" {
			valid = false
			errMsg = "Код не был прочитан"
		} else if len(code) != s.codeLength {
			valid = false
			errMsg = fmt.Sprintf("Код должен содержать %d символов, а в коде %d символов", s.codeLength, len(code))
		}

		// Проверяем на дубликаты
		if valid {
			if _, exists := s.codes[code]; exists {
				valid = false
				errMsg = "Дубликат кода"
			} else {
				// Добавляем код в список обработанных
				s.codes[code] = true
			}
		}

		// Обновляем статистику
		if valid {
			s.stats.ValidCodes++
		} else {
			s.stats.InvalidCodes++
		}

		s.stats.LastResult = valid
		s.stats.LastError = errMsg

		// Обновляем недавние коды (макс. 10)
		if len(s.stats.RecentCodes) >= 10 {
			// Сдвигаем список, убирая самый старый код
			s.stats.RecentCodes = append(s.stats.RecentCodes[1:], code)
			s.stats.RecentResults = append(s.stats.RecentResults[1:], valid)
		} else {
			// Добавляем код в список
			s.stats.RecentCodes = append(s.stats.RecentCodes, code)
			s.stats.RecentResults = append(s.stats.RecentResults, valid)
		}

		// Записываем валидный код в файл ролика
		if valid && s.rollManager != nil && s.rollManager.HasActiveRoll() {
			err := s.rollManager.WriteCode(code)
			if err != nil {
				// Обработка ошибки записи
				s.stats.LastError = "Ошибка записи в файл: " + err.Error()
			}
		}

		s.mu.Unlock()
	}
}
