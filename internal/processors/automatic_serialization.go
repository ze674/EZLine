package processors

import (
	"context"
	"sync"
)

type AutomaticSerializationProcessor struct {
	mu sync.Mutex
}

func NewAutomaticSerializationProcessor() *AutomaticSerializationProcessor {
	return &AutomaticSerializationProcessor{}
}

func (p *AutomaticSerializationProcessor) Start(ctx context.Context) error {
	return nil
}

func (p *AutomaticSerializationProcessor) Stop() error {
	return nil
}

func (p *AutomaticSerializationProcessor) IsRunning() bool {
	return false
}
