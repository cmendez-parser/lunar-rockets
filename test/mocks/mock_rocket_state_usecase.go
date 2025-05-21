package mocks

import (
	"context"
	"lunar-rockets/domain"

	"github.com/stretchr/testify/mock"
)

// MockRocketStateUsecase is a mock implementation of usecase.RocketStateUsecase
type MockRocketStateUsecase struct {
	mock.Mock
}

// UpdateRocketFromMessage calls the mocked implementation
func (m *MockRocketStateUsecase) UpdateRocketFromMessage(ctx context.Context, message *domain.RocketMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}
