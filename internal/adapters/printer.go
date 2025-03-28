package adapters

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Printer - структура для работы с принтером
type Printer struct {
	ip       string
	port     string
	conn     net.Conn
	mu       sync.Mutex
	isClosed bool
}

func NewPrinter(ip, port string) *Printer {
	return &Printer{ip: ip, port: port}
}

// Connect устанавливает соединение с принтером
func (p *Printer) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.isClosed {
		return nil // Соединение уже открыто
	}

	address := fmt.Sprintf("%s:%s", p.ip, p.port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к принтеру (%s): %v", address, err)
	}

	p.conn = conn
	p.isClosed = false
	return nil
}

// Send отправляет данные на принтер
func (p *Printer) Send(data string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn == nil || p.isClosed {
		return fmt.Errorf("соединение с принтером не установлено")
	}

	_, err := p.conn.Write([]byte(data))
	if err != nil {
		p.isClosed = true
		return fmt.Errorf("ошибка отправки данных: %v", err)
	}

	return nil
}

// Close закрывает соединение с принтером
func (p *Printer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil && !p.isClosed {
		err := p.conn.Close()
		if err != nil {
			return fmt.Errorf("ошибка закрытия соединения: %v", err)
		}
		p.isClosed = true
	}
	return nil
}
