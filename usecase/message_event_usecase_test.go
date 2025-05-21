package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"lunar-rockets/domain"
	"lunar-rockets/test/helper"
	"lunar-rockets/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessageEventUsecase_ProcessMessage(t *testing.T) {
	now := time.Now()
	testMessage := helper.CreateTestMessage("channel-1", domain.TypeRocketLaunched, 1, now)

	testCases := []struct {
		name              string
		message           *domain.RocketMessage
		lastMessageNumber int64
		messageRepoError  error
		stateUsecaseError error
		expectedError     string
		shouldCallState   bool
	}{
		{
			name:              "successful_first_message",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   true,
		},
		{
			name:              "successful_sequential_message",
			message:           helper.CreateTestMessage("channel-1", domain.TypeRocketSpeedIncreased, 2, now),
			lastMessageNumber: 1,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   true,
		},
		{
			name:              "duplicate_message",
			message:           testMessage,
			lastMessageNumber: 1,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   false,
		},
		{
			name:              "out_of_order_message",
			message:           helper.CreateTestMessage("channel-1", domain.TypeRocketSpeedIncreased, 2, now),
			lastMessageNumber: 0,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "message out of order: %!w(<nil>)",
			shouldCallState:   false,
		},
		{
			name:              "message_repo_error",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  errors.New("database error"),
			stateUsecaseError: nil,
			expectedError:     "failed to check if message was processed: database error",
			shouldCallState:   false,
		},
		{
			name:              "state_usecase_error",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  nil,
			stateUsecaseError: errors.New("state processing error"),
			expectedError:     "failed to execute rocket state usecase: state processing error",
			shouldCallState:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create mock repositories
			mockMessageRepo := &mocks.MockMessageRepository{
				FindLastMessageNumberFunc: func(ctx context.Context, channel string) (int64, error) {
					assert.Equal(t, tc.message.Metadata.Channel, channel)
					return tc.lastMessageNumber, tc.messageRepoError
				},
			}

			mockRocketRepo := &mocks.MockRocketRepository{}

			// Create mock rocket state usecase
			mockRocketStateUsecase := &mocks.MockRocketStateUsecase{}
			if tc.shouldCallState {
				if tc.stateUsecaseError == nil {
					mockRocketStateUsecase.On("UpdateRocketFromMessage", mock.Anything, tc.message).Return(nil)
				} else {
					mockRocketStateUsecase.On("UpdateRocketFromMessage", mock.Anything, tc.message).Return(tc.stateUsecaseError)
				}
			}

			// Create use case with mock dependencies
			useCase := NewMessageEventUsecase(mockRocketRepo, mockMessageRepo, mockRocketStateUsecase)

			// Execute the method
			err := useCase.ProcessMessage(context.Background(), tc.message)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRocketStateUsecase.AssertExpectations(t)
		})
	}
}
