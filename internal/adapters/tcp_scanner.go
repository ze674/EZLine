package adapters

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	connectTimeout = 5 * time.Second
	readTimeout    = 200 * time.Millisecond
	writeTimeout   = 200 * time.Millisecond
)

const suffix = '\n'

type Scanner struct {
	client      net.Conn
	address     string
	port        string
	scanCommand string
	reader      *bufio.Reader
}

// NewScanner создает новый экземпляр Scanner без установления соединения
func NewScanner(address, scanCommand string) *Scanner {
	return &Scanner{
		address:     address,
		scanCommand: scanCommand,
	}
}

// Connect устанавливает соединение
func (s *Scanner) Connect() error {
	op := "scanner.tcp.Connect"
	fmt.Println(op)
	client, err := net.DialTimeout("tcp", s.address, connectTimeout) // Устанавливаем соединение
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.client = client
	s.reader = bufio.NewReader(s.client) // Подготавливаем буфер для чтения данных

	return nil
}

// Close закрывает соединение
func (s *Scanner) Close() error {
	op := "scanner.tcp.Close"
	fmt.Println(op)

	var err error
	if s.client != nil {
		err = s.client.Close() // даже если ошибка, всё равно nil-им клиент
		s.client = nil
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Scan выполняет цикл сканирования
func (s *Scanner) Scan() (string, error) {
	op := "scanner.tcp.Scan"
	if s.client == nil {
		return "", fmt.Errorf("%s: scanner not connected", op)
	}
	fmt.Println(op)

	if err := s.SendCommand(s.scanCommand); err != nil { // Отправляем команду сканирования
		return "", fmt.Errorf("%s: %w", op, err)
	}

	response, err := s.ReadResponse() // Читаем ответ
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return response, nil

}

func (s *Scanner) SendCommand(command string) error {
	op := "scanner.tcp.SendCommand"

	if err := s.client.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil { // Устанавливаем таймаут для записи
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err := s.client.Write([]byte(command)) // Отправляем команду
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// ReadResponse читает ответ
func (s *Scanner) ReadResponse() (string, error) {
	op := "scanner.tcp.ReadResponse"

	if err := s.client.SetReadDeadline(time.Now().Add(readTimeout)); err != nil { // Устанавливаем таймаут для чтения
		return "", fmt.Errorf("%s: %w", op, err)
	}

	data, err := s.reader.ReadBytes(suffix) // Читаем ответ
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	response := string(data)
	response = strings.TrimSpace(response)
	return response, nil

}
