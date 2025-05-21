package mocks

import (
	"context"

	"lunar-rockets/domain"

	"github.com/stretchr/testify/mock"
)

type MockRocketMessageUsecase struct {
	mock.Mock
}

func (m *MockRocketMessageUsecase) ProcessMessage(ctx context.Context, message *domain.RocketMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}
