package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMessageRepository_MarkAsProcessed(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)

	testCases := []struct {
		name          string
		channel       string
		messageNumber int64
		expectedError string
	}{
		{
			name:          "successful_mark",
			channel:       "channel-1",
			messageNumber: 1,
			expectedError: "",
		},
		{
			name:          "database_error",
			channel:       "channel-1",
			messageNumber: 1,
			expectedError: "failed to mark message as processed: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectExec("INSERT INTO processed_messages \\(channel, message_number, processed_at\\)").
					WithArgs(tc.channel, tc.messageNumber).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO processed_messages \\(channel, message_number, processed_at\\)").
					WithArgs(tc.channel, tc.messageNumber).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			err := repo.MarkAsProcessed(context.Background(), tc.channel, tc.messageNumber)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessageRepository_FindLastMessageNumber(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)

	testCases := []struct {
		name           string
		channel        string
		mockRows       *sqlmock.Rows
		expectedNumber int64
		expectedError  string
	}{
		{
			name:    "has_messages",
			channel: "channel-1",
			mockRows: sqlmock.NewRows([]string{"MAX(message_number)"}).
				AddRow(5),
			expectedNumber: 5,
			expectedError:  "",
		},
		{
			name:           "no_messages",
			channel:        "channel-1",
			mockRows:       sqlmock.NewRows([]string{"MAX(message_number)"}),
			expectedNumber: 0,
			expectedError:  "",
		},
		{
			name:           "database_error",
			channel:        "channel-1",
			mockRows:       nil,
			expectedNumber: 0,
			expectedError:  "failed to find last message number: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectQuery("SELECT MAX\\(message_number\\) FROM processed_messages WHERE channel = \\?").
					WithArgs(tc.channel).
					WillReturnRows(tc.mockRows)
			} else {
				mock.ExpectQuery("SELECT MAX\\(message_number\\) FROM processed_messages WHERE channel = \\?").
					WithArgs(tc.channel).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			number, err := repo.FindLastMessageNumber(context.Background(), tc.channel)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Equal(t, int64(0), number)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedNumber, number)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
