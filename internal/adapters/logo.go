package adapters

import (
	"context"
	"fmt"
	"github.com/goburrow/modbus"
	"time"
)

const (
	on  uint16 = 0xFF00
	off uint16 = 0x0000
)

type Plc struct {
	address        string
	port           string
	client         modbus.Client
	handler        *modbus.TCPClientHandler
	sensorScanTime time.Duration
}

// New создает новый экземпляр Plc без установления соединения
func NewLogo(address string, sensorScanTime int) *Plc {
	return &Plc{
		address:        address,
		sensorScanTime: time.Duration(sensorScanTime),
	}
}

// Connect устанавливает соединение
func (p *Plc) Connect() error {
	op := "plc.modbus.Connect"

	p.handler = modbus.NewTCPClientHandler(p.address)
	p.handler.Timeout = 5 * time.Second

	if err := p.handler.Connect(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	client := modbus.NewClient(p.handler)
	p.client = client
	return nil
}

// Close закрывает соединение
func (p *Plc) Close() error {
	op := "plc.modbus.Close"

	if p.handler != nil {
		if err := p.handler.Close(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		p.handler = nil
	}
	return nil
}

// HandleProductSignal считывает состояние регистра датчика продукта и возвращает канал если состояние TRUE
func (p *Plc) HandleProductSignal(ctx context.Context, productSensorRegister uint16) (<-chan struct{}, error) {
	op := "plc.modbus.HandleProductSignal"
	ch := make(chan struct{}, 5) //Создаем канал

	lastState := false // Предыдущее состояние

	go func() {
		// Закрываем канал
		defer close(ch)
		for {
			select {
			case <-time.After(p.sensorScanTime): // Ждем p.sensorScanTime

				res, err := p.client.ReadCoils(productSensorRegister, 1) // Считываем состояние регистра
				if err != nil {
					fmt.Printf("%s: %s\n", op, err)
					continue
				}

				if len(res) == 0 { // Проверка на пустой ответ
					fmt.Printf("%s: %s\n", op, "empty response")
					continue
				}

				firstByte := res[0] // Получаем первый байт

				currentState := (firstByte & 0x01) == 0x01 // Проверяем байт на состояние

				if lastState == false && currentState == true { // Если был выключен и сейчас включен, то выключаем реле
					select {
					case ch <- struct{}{}: // Записываем в канал
					default:
						fmt.Printf("%s: %s\n", op, "channel is full") // Если канал заполнен, то выводим ошибку
					}
				}

				lastState = currentState // Запоминаем состояние

			case <-ctx.Done():
				return
			}

		}
	}()
	return ch, nil
}

// RejectorOn включает реле отбраковки
func (p *Plc) RejectorOn(rejectorRegister uint16) error {
	op := "plc.modbus.RejectorOn"

	_, err := p.client.WriteSingleCoil(rejectorRegister, on) // Включаем реле
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// RejectorOff выключает реле отбраковки
func (p *Plc) RejectorOff(rejectorRegister uint16) error {
	op := "plc.modbus.RejectorOff"

	_, err := p.client.WriteSingleCoil(rejectorRegister, off) // Отключаем реле
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
