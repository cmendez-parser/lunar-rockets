package mocks

import (
	"context"
	"lunar-rockets/domain"
)

// MockRocketRepository is a mock implementation of domain.RocketRepository
type MockRocketRepository struct {
	GetByChannelFunc func(ctx context.Context, channel string) (*domain.Rocket, error)
	GetAllFunc       func(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error)
	SaveFunc         func(ctx context.Context, rocket *domain.Rocket) error
	UpdateFunc       func(ctx context.Context, rocket *domain.Rocket) error
	DeleteFunc       func(ctx context.Context, channel string) error
	BeginTxFunc      func(ctx context.Context) (domain.Transaction, error)
}

// Ensure MockRocketRepository implements domain.RocketRepository
var _ domain.RocketRepository = (*MockRocketRepository)(nil)

// GetByChannel calls the mocked implementation
func (m *MockRocketRepository) GetByChannel(ctx context.Context, channel string) (*domain.Rocket, error) {
	return m.GetByChannelFunc(ctx, channel)
}

// GetAll calls the mocked implementation
func (m *MockRocketRepository) GetAll(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error) {
	return m.GetAllFunc(ctx, sortBy, order)
}

// Save calls the mocked implementation
func (m *MockRocketRepository) Save(ctx context.Context, rocket *domain.Rocket) error {
	return m.SaveFunc(ctx, rocket)
}

// Update calls the mocked implementation
func (m *MockRocketRepository) Update(ctx context.Context, rocket *domain.Rocket) error {
	return m.UpdateFunc(ctx, rocket)
}

// Delete calls the mocked implementation
func (m *MockRocketRepository) Delete(ctx context.Context, channel string) error {
	return m.DeleteFunc(ctx, channel)
}

// BeginTx calls the mocked implementation
func (m *MockRocketRepository) BeginTx(ctx context.Context) (domain.Transaction, error) {
	return m.BeginTxFunc(ctx)
}

// MockTransaction is a mock implementation of domain.Transaction
type MockTransaction struct {
	CommitFunc   func() error
	RollbackFunc func() error
}

// Ensure MockTransaction implements domain.Transaction
var _ domain.Transaction = (*MockTransaction)(nil)

// Commit calls the mocked implementation
func (m *MockTransaction) Commit() error {
	return m.CommitFunc()
}

// Rollback calls the mocked implementation
func (m *MockTransaction) Rollback() error {
	return m.RollbackFunc()
}
