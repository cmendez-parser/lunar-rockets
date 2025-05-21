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
)

func TestRocketStateUsecase_UpdateRocketFromMessage(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name                string
		message             *domain.RocketMessage
		existingRocket      *domain.Rocket
		rocketRepoError     error
		messageRepoError    error
		expectedError       string
		expectedRocketState *domain.Rocket
		ignoreLastUpdated   bool // Flag to ignore LastUpdated field comparison
		ignoreRocketState   bool // Flag to ignore rocket state comparison
	}{
		{
			name: "successful_rocket_launch",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketLaunched,
					MessageNumber: 1,
					MessageTime:   now,
				},
				Message: domain.RocketLaunchedMessage{
					Type:        "Falcon-9",
					LaunchSpeed: 1000,
					Mission:     "ARTEMIS",
				},
			},
			existingRocket:   nil,
			rocketRepoError:  nil,
			messageRepoError: nil,
			expectedError:    "",
			expectedRocketState: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1000,
				Mission:     "ARTEMIS",
				LaunchTime:  now,
				Status:      domain.RocketStatusLaunched,
				LastMessage: 1,
			},
			ignoreLastUpdated: true,
		},
		{
			name: "successful_speed_increase",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketSpeedIncreased,
					MessageNumber: 2,
					MessageTime:   now,
				},
				Message: domain.RocketSpeedIncreasedMessage{
					By: 500,
				},
			},
			existingRocket:   helper.CreateTestRocket("channel-1", "Falcon-9", "ARTEMIS", domain.RocketStatusLaunched, 1000, now.Add(-1*time.Hour)),
			rocketRepoError:  nil,
			messageRepoError: nil,
			expectedError:    "",
			expectedRocketState: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1500,
				Mission:     "ARTEMIS",
				LaunchTime:  now.Add(-1 * time.Hour),
				Status:      domain.RocketStatusLaunched,
				LastMessage: 2,
			},
			ignoreLastUpdated: true,
		},
		{
			name: "successful_speed_decrease",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketSpeedDecreased,
					MessageNumber: 3,
					MessageTime:   now,
				},
				Message: domain.RocketSpeedDecreasedMessage{
					By: 300,
				},
			},
			existingRocket:   helper.CreateTestRocket("channel-1", "Falcon-9", "ARTEMIS", domain.RocketStatusLaunched, 1500, now.Add(-1*time.Hour)),
			rocketRepoError:  nil,
			messageRepoError: nil,
			expectedError:    "",
			expectedRocketState: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1200,
				Mission:     "ARTEMIS",
				LaunchTime:  now.Add(-1 * time.Hour),
				Status:      domain.RocketStatusLaunched,
				LastMessage: 3,
			},
			ignoreLastUpdated: true,
		},
		{
			name: "successful_mission_change",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketMissionChanged,
					MessageNumber: 4,
					MessageTime:   now,
				},
				Message: domain.RocketMissionChangedMessage{
					NewMission: "MARS",
				},
			},
			existingRocket:   helper.CreateTestRocket("channel-1", "Falcon-9", "ARTEMIS", domain.RocketStatusLaunched, 1200, now.Add(-1*time.Hour)),
			rocketRepoError:  nil,
			messageRepoError: nil,
			expectedError:    "",
			expectedRocketState: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1200,
				Mission:     "MARS",
				LaunchTime:  now.Add(-1 * time.Hour),
				Status:      domain.RocketStatusLaunched,
				LastMessage: 4,
			},
			ignoreLastUpdated: true,
		},
		{
			name: "successful_rocket_explosion",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketExploded,
					MessageNumber: 5,
					MessageTime:   now,
				},
				Message: domain.RocketExplodedMessage{
					Reason: "PRESSURE_FAILURE",
				},
			},
			existingRocket:   helper.CreateTestRocket("channel-1", "Falcon-9", "MARS", domain.RocketStatusLaunched, 1200, now.Add(-1*time.Hour)),
			rocketRepoError:  nil,
			messageRepoError: nil,
			expectedError:    "",
			expectedRocketState: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1200,
				Mission:     "MARS",
				LaunchTime:  now.Add(-1 * time.Hour),
				Status:      domain.RocketStatusExploded,
				Reason:      "PRESSURE_FAILURE",
				ExplodedAt:  &now,
				LastMessage: 5,
			},
			ignoreLastUpdated: true,
		},
		{
			name: "rocket_not_found_for_update",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketSpeedIncreased,
					MessageNumber: 2,
					MessageTime:   now,
				},
				Message: domain.RocketSpeedIncreasedMessage{
					By: 500,
				},
			},
			existingRocket:      nil,
			rocketRepoError:     nil,
			messageRepoError:    nil,
			expectedError:       "failed to process message: rocket not found: channel-1",
			expectedRocketState: nil,
		},
		{
			name: "rocket_repo_error",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketLaunched,
					MessageNumber: 1,
					MessageTime:   now,
				},
				Message: domain.RocketLaunchedMessage{
					Type:        "Falcon-9",
					LaunchSpeed: 1000,
					Mission:     "ARTEMIS",
				},
			},
			existingRocket:      nil,
			rocketRepoError:     errors.New("database error"),
			messageRepoError:    nil,
			expectedError:       "failed to process message: database error",
			expectedRocketState: nil,
		},
		{
			name: "message_repo_error",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   domain.TypeRocketLaunched,
					MessageNumber: 1,
					MessageTime:   now,
				},
				Message: domain.RocketLaunchedMessage{
					Type:        "Falcon-9",
					LaunchSpeed: 1000,
					Mission:     "ARTEMIS",
				},
			},
			existingRocket:      nil,
			rocketRepoError:     nil,
			messageRepoError:    errors.New("database error"),
			expectedError:       "failed to mark message as processed: database error",
			expectedRocketState: nil,
			ignoreRocketState:   true, // We don't care about the rocket state in this case
		},
		{
			name: "unknown_message_type",
			message: &domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageType:   "UnknownType",
					MessageNumber: 1,
					MessageTime:   now,
				},
				Message: nil,
			},
			existingRocket:      nil,
			rocketRepoError:     nil,
			messageRepoError:    nil,
			expectedError:       "failed to process message: unknown message type: UnknownType",
			expectedRocketState: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create mock repositories
			mockRocketRepo := &mocks.MockRocketRepository{
				GetByChannelFunc: func(ctx context.Context, channel string) (*domain.Rocket, error) {
					assert.Equal(t, tc.message.Metadata.Channel, channel)
					return tc.existingRocket, tc.rocketRepoError
				},
				SaveFunc: func(ctx context.Context, rocket *domain.Rocket) error {
					if tc.ignoreRocketState {
						return tc.rocketRepoError
					}
					if tc.ignoreLastUpdated {
						// Create a copy of the expected rocket with the actual LastUpdated time
						expectedRocket := *tc.expectedRocketState
						expectedRocket.LastUpdated = rocket.LastUpdated
						assert.Equal(t, &expectedRocket, rocket)
					} else {
						assert.Equal(t, tc.expectedRocketState, rocket)
					}
					return tc.rocketRepoError
				},
				UpdateFunc: func(ctx context.Context, rocket *domain.Rocket) error {
					if tc.ignoreRocketState {
						return tc.rocketRepoError
					}
					if tc.ignoreLastUpdated {
						// Create a copy of the expected rocket with the actual LastUpdated time
						expectedRocket := *tc.expectedRocketState
						expectedRocket.LastUpdated = rocket.LastUpdated
						assert.Equal(t, &expectedRocket, rocket)
					} else {
						assert.Equal(t, tc.expectedRocketState, rocket)
					}
					return tc.rocketRepoError
				},
				BeginTxFunc: func(ctx context.Context) (domain.Transaction, error) {
					return &mocks.MockTransaction{
						CommitFunc: func() error {
							return nil
						},
						RollbackFunc: func() error {
							return nil
						},
					}, nil
				},
			}

			mockMessageRepo := &mocks.MockMessageRepository{
				MarkAsProcessedFunc: func(ctx context.Context, channel string, messageNumber int64) error {
					assert.Equal(t, tc.message.Metadata.Channel, channel)
					assert.Equal(t, tc.message.Metadata.MessageNumber, messageNumber)
					return tc.messageRepoError
				},
			}

			// Create use case with mock dependencies
			useCase := NewRocketStateUsecase(mockRocketRepo, mockMessageRepo)

			// Execute the method
			err := useCase.UpdateRocketFromMessage(context.Background(), tc.message)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
