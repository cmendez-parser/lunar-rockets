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

func TestRocketMessageUsecase_ProcessMessage(t *testing.T) {
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
		shouldBuffer      bool
	}{
		{
			name:              "successful_first_message",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   true,
			shouldBuffer:      false,
		},
		{
			name:              "successful_sequential_message",
			message:           helper.CreateTestMessage("channel-1", domain.TypeRocketSpeedIncreased, 2, now),
			lastMessageNumber: 1,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   true,
			shouldBuffer:      false,
		},
		{
			name:              "duplicate_message",
			message:           testMessage,
			lastMessageNumber: 1,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   false,
			shouldBuffer:      false,
		},
		{
			name:              "out_of_order_message",
			message:           helper.CreateTestMessage("channel-1", domain.TypeRocketSpeedIncreased, 3, now),
			lastMessageNumber: 1,
			messageRepoError:  nil,
			stateUsecaseError: nil,
			expectedError:     "",
			shouldCallState:   false,
			shouldBuffer:      true,
		},
		{
			name:              "message_repo_error",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  errors.New("database error"),
			stateUsecaseError: nil,
			expectedError:     "failed to check if message was processed: database error",
			shouldCallState:   false,
			shouldBuffer:      false,
		},
		{
			name:              "state_usecase_error",
			message:           testMessage,
			lastMessageNumber: 0,
			messageRepoError:  nil,
			stateUsecaseError: errors.New("state processing error"),
			expectedError:     "failed to execute rocket state usecase: state processing error",
			shouldCallState:   true,
			shouldBuffer:      false,
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
			useCase := NewRocketMessageUsecase(mockRocketRepo, mockMessageRepo, mockRocketStateUsecase)

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

			// Verify buffer state
			if tc.shouldBuffer {
				assert.Contains(t, useCase.(*rocketMessageUsecase).messageBuffer[tc.message.Metadata.Channel], tc.message.Metadata.MessageNumber)
			} else {
				assert.NotContains(t, useCase.(*rocketMessageUsecase).messageBuffer[tc.message.Metadata.Channel], tc.message.Metadata.MessageNumber)
			}
		})
	}
}

func TestRocketMessageUsecase_ProcessBufferedMessages(t *testing.T) {
	now := time.Now()
	channel := "test-channel"

	// Create a sequence of messages
	messages := []*domain.RocketMessage{
		helper.CreateTestMessage(channel, domain.TypeRocketLaunched, 1, now),
		helper.CreateTestMessage(channel, domain.TypeRocketSpeedIncreased, 2, now),
		helper.CreateTestMessage(channel, domain.TypeRocketSpeedIncreased, 3, now),
	}

	testCases := []struct {
		name                string
		lastProcessedNumber int64
		stateUsecaseError   error
		expectedError       string
		expectedProcessed   []int64 // Message numbers that should be processed
	}{
		{
			name:                "process_all_buffered_messages",
			lastProcessedNumber: 0,
			stateUsecaseError:   nil,
			expectedError:       "",
			expectedProcessed:   []int64{1, 2, 3},
		},
		{
			name:                "process_partial_buffered_messages",
			lastProcessedNumber: 1,
			stateUsecaseError:   nil,
			expectedError:       "",
			expectedProcessed:   []int64{2, 3}, // Should process both 2 and 3 since they're consecutive
		},
		{
			name:                "error_processing_buffered_message",
			lastProcessedNumber: 0,
			stateUsecaseError:   errors.New("state processing error"),
			expectedError:       "failed to process buffered message 1: state processing error",
			expectedProcessed:   []int64{1},
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create mock repositories
			mockMessageRepo := &mocks.MockMessageRepository{}
			mockRocketRepo := &mocks.MockRocketRepository{}

			// Create mock rocket state usecase
			mockRocketStateUsecase := &mocks.MockRocketStateUsecase{}

			// Set up mock expectations for each message that should be processed
			for _, msgNum := range tc.expectedProcessed {
				msg := messages[msgNum-1]
				if tc.stateUsecaseError != nil && msgNum == tc.expectedProcessed[0] {
					mockRocketStateUsecase.On("UpdateRocketFromMessage", mock.Anything, msg).Return(tc.stateUsecaseError)
				} else {
					mockRocketStateUsecase.On("UpdateRocketFromMessage", mock.Anything, msg).Return(nil)
				}
			}

			// Create use case with mock dependencies
			useCase := NewRocketMessageUsecase(mockRocketRepo, mockMessageRepo, mockRocketStateUsecase)

			// Add messages to buffer
			for _, msg := range messages {
				useCase.(*rocketMessageUsecase).addToBuffer(msg)
			}

			// Process buffered messages
			err := useCase.(*rocketMessageUsecase).processBufferedMessages(context.Background(), channel, tc.lastProcessedNumber)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRocketStateUsecase.AssertExpectations(t)

			// Verify buffer state after processing
			if tc.stateUsecaseError == nil {
				// Verify that processed messages were removed from buffer
				for _, msgNum := range tc.expectedProcessed {
					assert.NotContains(t, useCase.(*rocketMessageUsecase).messageBuffer[channel], msgNum, "Processed message should be removed from buffer")
				}

				// Verify that unprocessed messages are still in the buffer
				for _, msg := range messages {
					if !contains(tc.expectedProcessed, msg.Metadata.MessageNumber) {
						assert.Contains(t, useCase.(*rocketMessageUsecase).messageBuffer[channel], msg.Metadata.MessageNumber, "Unprocessed message should remain in buffer")
					}
				}
			}
		})
	}
}

// Helper function to check if a slice contains a value
func contains(slice []int64, value int64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
