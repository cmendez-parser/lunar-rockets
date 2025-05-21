package controller

import (
	"bytes"
	"encoding/json"
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

func TestNewMessageController(t *testing.T) {
	mockUsecase := &mocks.MockMessageEventUsecase{}
	controller := NewMessageController(mockUsecase)

	assert.NotNil(t, controller)
	assert.Equal(t, mockUsecase, controller.messageEventUsecase)
}

func TestMessageController_ReceiveMessage(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		body           interface{}
		setupMock      func(*mocks.MockMessageEventUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "valid_message",
			method: http.MethodPost,
			body: domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageNumber: 1,
					MessageTime:   time.Now(),
					MessageType:   domain.TypeRocketLaunched,
				},
				Message: domain.RocketLaunchedMessage{
					Type:        "Falcon-9",
					LaunchSpeed: 1000,
					Mission:     "ARTEMIS",
				},
			},
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				m.On("ProcessMessage", mock.Anything, mock.AnythingOfType("*domain.RocketMessage")).
					Return(nil)
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"status":"accepted"}`,
		},
		{
			name:   "invalid_method",
			method: http.MethodGet,
			body:   nil,
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
		{
			name:   "invalid_json",
			method: http.MethodPost,
			body:   "invalid json",
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid message format\n",
		},
		{
			name:   "missing_channel",
			method: http.MethodPost,
			body: domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					MessageNumber: 1,
					MessageTime:   time.Now(),
					MessageType:   domain.TypeRocketLaunched,
				},
			},
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing channel ID\n",
		},
		{
			name:   "missing_message_type",
			method: http.MethodPost,
			body: domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageNumber: 1,
					MessageTime:   time.Now(),
				},
			},
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing message type\n",
		},
		{
			name:   "database_error",
			method: http.MethodPost,
			body: domain.RocketMessage{
				Metadata: domain.MessageMetadata{
					Channel:       "channel-1",
					MessageNumber: 1,
					MessageTime:   time.Now(),
					MessageType:   domain.TypeRocketLaunched,
				},
				Message: domain.RocketLaunchedMessage{
					Type:        "Falcon-9",
					LaunchSpeed: 1000,
					Mission:     "ARTEMIS",
				},
			},
			setupMock: func(m *mocks.MockMessageEventUsecase) {
				m.On("ProcessMessage", mock.Anything, mock.AnythingOfType("*domain.RocketMessage")).
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to process message\n",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock for each test case
			mockUsecase := &mocks.MockMessageEventUsecase{}
			controller := NewMessageController(mockUsecase)

			// Setup mock
			tc.setupMock(mockUsecase)

			// Create request
			var body []byte
			var err error
			if tc.body != nil {
				body, err = json.Marshal(tc.body)
				assert.NoError(t, err)
			} else {
				body = []byte("invalid json")
			}

			req := httptest.NewRequest(tc.method, "/messages", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			// Execute request
			controller.ReceiveMessage(w, req)

			// Check response
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

			// Verify mock expectations
			mockUsecase.AssertExpectations(t)
		})
	}
}
