package mocks

import (
	"context"
	"lunar-rockets/domain"
)

// MockMessageRepository is a mock implementation of domain.MessageRepository
type MockMessageRepository struct {
	MarkAsProcessedFunc       func(ctx context.Context, channel string, messageNumber int64) error
	FindLastMessageNumberFunc func(ctx context.Context, channel string) (int64, error)
}

// Ensure MockMessageRepository implements domain.MessageRepository
var _ domain.MessageRepository = (*MockMessageRepository)(nil)

// MarkAsProcessed calls the mocked implementation
func (m *MockMessageRepository) MarkAsProcessed(ctx context.Context, channel string, messageNumber int64) error {
	return m.MarkAsProcessedFunc(ctx, channel, messageNumber)
}

// FindLastMessageNumber calls the mocked implementation
func (m *MockMessageRepository) FindLastMessageNumber(ctx context.Context, channel string) (int64, error) {
	return m.FindLastMessageNumberFunc(ctx, channel)
}
