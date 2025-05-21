package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lunar-rockets/domain"
	"lunar-rockets/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var fixedTime = time.Date(2024, 3, 21, 0, 0, 0, 0, time.UTC)

func TestNewRocketController(t *testing.T) {
	mockUsecase := &mocks.MockRocketUseCase{}
	controller := NewRocketController(mockUsecase)

	assert.NotNil(t, controller)
	assert.Equal(t, mockUsecase, controller.rocketUseCase)
}

func TestRocketController_GetRocket(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		channelID      string
		setupMock      func(*mocks.MockRocketUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "valid_rocket",
			method:    http.MethodGet,
			channelID: "channel-1",
			setupMock: func(m *mocks.MockRocketUseCase) {
				m.On("GetRocket", mock.Anything, "channel-1").
					Return(&domain.Rocket{
						Channel:     "channel-1",
						Type:        "Falcon-9",
						Speed:       1000,
						Mission:     "ARTEMIS",
						Status:      domain.RocketStatusLaunched,
						LaunchTime:  fixedTime,
						LastUpdated: fixedTime,
						LastMessage: 1,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"channel":"channel-1","type":"Falcon-9","speed":1000,"mission":"ARTEMIS","launchTime":"2024-03-21T00:00:00Z","status":"Launched","lastUpdated":"2024-03-21T00:00:00Z","lastMessage":1}` + "\n",
		},
		{
			name:      "invalid_method",
			method:    http.MethodPost,
			channelID: "channel-1",
			setupMock: func(m *mocks.MockRocketUseCase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
		{
			name:      "missing_channel",
			method:    http.MethodGet,
			channelID: "",
			setupMock: func(m *mocks.MockRocketUseCase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing channel ID\n",
		},
		{
			name:      "rocket_not_found",
			method:    http.MethodGet,
			channelID: "channel-1",
			setupMock: func(m *mocks.MockRocketUseCase) {
				m.On("GetRocket", mock.Anything, "channel-1").
					Return(nil, domain.ErrRocketNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Rocket not found\n",
		},
		{
			name:      "database_error",
			method:    http.MethodGet,
			channelID: "channel-1",
			setupMock: func(m *mocks.MockRocketUseCase) {
				m.On("GetRocket", mock.Anything, "channel-1").
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get rocket state\n",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock for each test case
			mockUsecase := &mocks.MockRocketUseCase{}
			controller := NewRocketController(mockUsecase)

			// Setup mock
			tc.setupMock(mockUsecase)

			// Create request
			req := httptest.NewRequest(tc.method, "/rockets/"+tc.channelID, nil)
			w := httptest.NewRecorder()

			// Execute request
			controller.GetRocket(w, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

			// Verify mock expectations
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestRocketController_ListRockets(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		query          string
		setupMock      func(*mocks.MockRocketUseCase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "valid_list",
			method: http.MethodGet,
			query:  "?sort=speed&order=desc",
			setupMock: func(m *mocks.MockRocketUseCase) {
				m.On("ListRockets", mock.Anything, "speed", "desc").
					Return([]*domain.Rocket{
						{
							Channel:     "channel-1",
							Type:        "Falcon-9",
							Speed:       1000,
							Mission:     "ARTEMIS",
							Status:      domain.RocketStatusLaunched,
							LaunchTime:  fixedTime,
							LastUpdated: fixedTime,
							LastMessage: 1,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"channel":"channel-1","type":"Falcon-9","speed":1000,"mission":"ARTEMIS","launchTime":"2024-03-21T00:00:00Z","status":"Launched","lastUpdated":"2024-03-21T00:00:00Z","lastMessage":1}]` + "\n",
		},
		{
			name:   "invalid_method",
			method: http.MethodPost,
			query:  "",
			setupMock: func(m *mocks.MockRocketUseCase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
		{
			name:   "database_error",
			method: http.MethodGet,
			query:  "?sort=speed&order=desc",
			setupMock: func(m *mocks.MockRocketUseCase) {
				m.On("ListRockets", mock.Anything, "speed", "desc").
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get rockets\n",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock for each test case
			mockUsecase := &mocks.MockRocketUseCase{}
			controller := NewRocketController(mockUsecase)

			// Setup mock
			tc.setupMock(mockUsecase)

			// Create request
			req := httptest.NewRequest(tc.method, "/rockets"+tc.query, nil)
			w := httptest.NewRecorder()

			// Execute request
			controller.ListRockets(w, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

			// Verify mock expectations
			mockUsecase.AssertExpectations(t)
		})
	}
}
