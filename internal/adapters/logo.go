package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/goburrow/modbus"
)

const (
	modbusOn  uint16 = 0xFF00
	modbusOff uint16 = 0x0000
)

// ModbusPLC реализует PLC-интерфейс через Modbus TCP
type ModbusPLC struct {
	address        string
	port           string
	timeout        time.Duration
	sensorScanTime time.Duration

	productSensorRegister uint16
	rejectorRegister      uint16

	client  modbus.Client
	handler *modbus.TCPClientHandler

	bufferSize int
}

// NewModbusPLC возвращает адаптер для PLC через Modbus TCP
func NewModbusPLC(address string, timeout, scanInterval time.Duration, productSensorRegister, rejectorRegister uint16, bufferSize int) *ModbusPLC {
	return &ModbusPLC{
		address:               address,
		timeout:               timeout,
		sensorScanTime:        scanInterval,
		productSensorRegister: productSensorRegister,
		rejectorRegister:      rejectorRegister,
		bufferSize:            bufferSize,
	}
}

// Connect устанавливает соединение с PLC
func (p *ModbusPLC) Connect() error {
	op := "plc.modbus.Connect"

	p.handler = modbus.NewTCPClientHandler(p.address)
	p.handler.Timeout = p.timeout

	if err := p.handler.Connect(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.client = modbus.NewClient(p.handler)
	return nil
}

// Close закрывает соединение с PLC
func (p *ModbusPLC) Close() error {
	op := "plc.modbus.Close"

	var err error
	if p.handler != nil {
		err = p.handler.Close()
		p.handler = nil
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// HandleProductSignal запускает мониторинг регистра и возвращает канал,
// в который отправляется сигнал при изменении с 0 на 1 (фронт)
func (p *ModbusPLC) HandleProductSignal(ctx context.Context) (<-chan struct{}, error) {
	op := "plc.modbus.HandleProductSignal"

	ch := make(chan struct{}, p.bufferSize)
	lastState := false

	go func() {
		defer close(ch)

		for {
			select {
			case <-time.After(p.sensorScanTime):
				res, err := p.client.ReadCoils(p.productSensorRegister, 1)
				if err != nil {
					fmt.Printf("%s: %s\n", op, err) // заменим на logger если появится
					continue
				}
				if len(res) == 0 {
					fmt.Printf("%s: %s\n", op, "empty response")
					continue
				}

				currentState := (res[0] & 0x01) == 0x01

				if !lastState && currentState {
					select {
					case ch <- struct{}{}:
					default:
						fmt.Printf("%s: channel full\n", op)
					}
				}
				lastState = currentState

			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// RejectorOn включает реле отбраковки
func (p *ModbusPLC) RejectorOn() error {
	op := "plc.modbus.RejectorOn"

	_, err := p.client.WriteSingleCoil(p.rejectorRegister, modbusOn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// RejectorOff выключает реле отбраковки
func (p *ModbusPLC) RejectorOff() error {
	op := "plc.modbus.RejectorOff"

	_, err := p.client.WriteSingleCoil(p.rejectorRegister, modbusOff)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
