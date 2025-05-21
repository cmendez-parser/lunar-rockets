package mocks

import (
	"context"

	"lunar-rockets/domain"

	"github.com/stretchr/testify/mock"
)

type MockMessageEventUsecase struct {
	mock.Mock
}

func (m *MockMessageEventUsecase) ProcessMessage(ctx context.Context, message *domain.RocketMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}
