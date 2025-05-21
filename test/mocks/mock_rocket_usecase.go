package mocks

import (
	"context"
	"lunar-rockets/domain"

	"github.com/stretchr/testify/mock"
)

// MockRocketUseCase is a mock implementation of usecase.RocketUseCase
type MockRocketUseCase struct {
	mock.Mock
}

func (m *MockRocketUseCase) GetRocket(ctx context.Context, channel string) (*domain.Rocket, error) {
	args := m.Called(ctx, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rocket), args.Error(1)
}

func (m *MockRocketUseCase) ListRockets(ctx context.Context, sortBy, order string) ([]*domain.Rocket, error) {
	args := m.Called(ctx, sortBy, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Rocket), args.Error(1)
}
