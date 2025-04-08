package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TimerTrigger реализует источник триггеров по таймеру
type TimerTrigger struct {
	interval time.Duration
	running  bool
	mu       sync.Mutex
	cancel   context.CancelFunc
	signal   chan struct{}
}

// NewTimerTrigger создает новый источник триггеров на основе таймера
func NewTimerTrigger(interval time.Duration) *TimerTrigger {
	return &TimerTrigger{
		interval: interval,
	}
}

// Start запускает таймер и возвращает канал с сигналами
func (t *TimerTrigger) WaitSignal(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return fmt.Errorf("таймер уже запущен")
	}

	triggerChan := make(chan struct{})

	// Создаем отдельный контекст для таймера, чтобы иметь возможность отменить его независимо
	triggerCtx, cancel := context.WithCancel(ctx)
	t.cancel = cancel

	// Запускаем таймер в отдельной горутине
	go func() {
		defer close(triggerChan)

		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Отправляем сигнал в канал
				select {
				case triggerChan <- struct{}{}:
					fmt.Println("Сигнал отправлен")
					// Успешно отправили сигнал
				default:
					// Канал заблокирован, пропускаем этот тик
				}
			case <-triggerCtx.Done():
				return
			}
		}
	}()

	t.signal = triggerChan
	t.running = true
	return nil
}

// Stop останавливает таймер
func (t *TimerTrigger) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return nil
	}

	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
	}

	t.running = false
	return nil
}

func (t *TimerTrigger) SignalChan() <-chan struct{} {
	return t.signal
}
