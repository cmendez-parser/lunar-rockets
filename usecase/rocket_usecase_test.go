package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"lunar-rockets/domain"
	"lunar-rockets/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRocketUseCase_GetRocket(t *testing.T) {
	fixedTime := time.Date(2025, 5, 20, 9, 39, 15, 253357000, time.Local)
	testCases := []struct {
		name           string
		channel        string
		rocket         *domain.Rocket
		repoError      error
		expectedError  string
		expectedRocket *domain.Rocket
	}{
		{
			name:    "successful_retrieval",
			channel: "test-channel",
			rocket: &domain.Rocket{
				Channel:     "test-channel",
				Type:        "test-type",
				Speed:       100,
				Mission:     "test-mission",
				LaunchTime:  fixedTime,
				Status:      domain.RocketStatusLaunched,
				LastUpdated: fixedTime,
			},
			repoError:     nil,
			expectedError: "",
			expectedRocket: &domain.Rocket{
				Channel:     "test-channel",
				Type:        "test-type",
				Speed:       100,
				Mission:     "test-mission",
				LaunchTime:  fixedTime,
				Status:      domain.RocketStatusLaunched,
				LastUpdated: fixedTime,
			},
		},
		{
			name:           "not_found",
			channel:        "nonexistent-channel",
			rocket:         nil,
			repoError:      nil,
			expectedError:  domain.ErrRocketNotFound.Error(),
			expectedRocket: nil,
		},
		{
			name:           "repository_error",
			channel:        "test-channel",
			rocket:         nil,
			repoError:      errors.New("database error"),
			expectedError:  "failed to get rocket: database error",
			expectedRocket: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &mocks.MockRocketRepository{
				GetByChannelFunc: func(ctx context.Context, channel string) (*domain.Rocket, error) {
					assert.Equal(t, tc.channel, channel)
					return tc.rocket, tc.repoError
				},
			}

			useCase := NewRocketUseCase(mockRepo)
			rocket, err := useCase.GetRocket(context.Background(), tc.channel)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRocket, rocket)
			}
		})
	}
}

func TestRocketUseCase_ListRockets(t *testing.T) {
	fixedTime := time.Date(2025, 5, 20, 9, 39, 15, 253357000, time.Local)
	testCases := []struct {
		name           string
		sortBy         string
		order          string
		rockets        []*domain.Rocket
		repoError      error
		expectedError  string
		expectedSortBy string
		expectedOrder  string
	}{
		{
			name:   "successful_list_default_sort",
			sortBy: "",
			order:  "",
			rockets: []*domain.Rocket{
				{
					Channel:     "test-channel-1",
					Type:        "test-type",
					Speed:       100,
					Mission:     "test-mission",
					LaunchTime:  fixedTime,
					Status:      domain.RocketStatusLaunched,
					LastUpdated: fixedTime,
				},
			},
			repoError:      nil,
			expectedError:  "",
			expectedSortBy: "launch_time",
			expectedOrder:  "DESC",
		},
		{
			name:   "successful_list_custom_sort",
			sortBy: "speed",
			order:  "ASC",
			rockets: []*domain.Rocket{
				{
					Channel:     "test-channel-1",
					Type:        "test-type",
					Speed:       100,
					Mission:     "test-mission",
					LaunchTime:  fixedTime,
					Status:      domain.RocketStatusLaunched,
					LastUpdated: fixedTime,
				},
			},
			repoError:      nil,
			expectedError:  "",
			expectedSortBy: "speed",
			expectedOrder:  "ASC",
		},
		{
			name:   "normalize_order",
			sortBy: "speed",
			order:  "invalid",
			rockets: []*domain.Rocket{
				{
					Channel:     "test-channel-1",
					Type:        "test-type",
					Speed:       100,
					Mission:     "test-mission",
					LaunchTime:  fixedTime,
					Status:      domain.RocketStatusLaunched,
					LastUpdated: fixedTime,
				},
			},
			repoError:      nil,
			expectedError:  "",
			expectedSortBy: "speed",
			expectedOrder:  "DESC",
		},
		{
			name:           "repository_error",
			sortBy:         "",
			order:          "",
			rockets:        nil,
			repoError:      errors.New("database error"),
			expectedError:  "failed to list rockets: database error",
			expectedSortBy: "launch_time",
			expectedOrder:  "DESC",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := &mocks.MockRocketRepository{
				GetAllFunc: func(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error) {
					assert.Equal(t, tc.expectedSortBy, sortBy)
					assert.Equal(t, tc.expectedOrder, order)
					return tc.rockets, tc.repoError
				},
			}

			useCase := NewRocketUseCase(mockRepo)
			rockets, err := useCase.ListRockets(context.Background(), tc.sortBy, tc.order)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.rockets, rockets)
			}
		})
	}
}
