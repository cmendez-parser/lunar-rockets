package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"lunar-rockets/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRocketRepository_GetByChannel(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	testCases := []struct {
		name           string
		channel        string
		mockRows       *sqlmock.Rows
		expectedRocket *domain.Rocket
		expectedError  string
	}{
		{
			name:    "successful_get",
			channel: "channel-1",
			mockRows: sqlmock.NewRows([]string{
				"channel", "type", "speed", "mission", "launch_time", "status",
				"exploded_at", "reason", "last_updated",
			}).AddRow(
				"channel-1", "Falcon-9", 1000, "ARTEMIS",
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				domain.RocketStatusLaunched,
				nil, nil,
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			),
			expectedRocket: &domain.Rocket{
				Channel:     "channel-1",
				Type:        "Falcon-9",
				Speed:       1000,
				Mission:     "ARTEMIS",
				LaunchTime:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:      domain.RocketStatusLaunched,
				LastUpdated: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: "",
		},
		{
			name:           "not_found",
			channel:        "channel-2",
			mockRows:       sqlmock.NewRows([]string{}),
			expectedRocket: nil,
			expectedError:  "",
		},
		{
			name:           "database_error",
			channel:        "channel-1",
			mockRows:       nil,
			expectedRocket: nil,
			expectedError:  "failed to get rocket: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectQuery("SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated FROM rockets WHERE channel = ?").
					WithArgs(tc.channel).
					WillReturnRows(tc.mockRows)
			} else {
				mock.ExpectQuery("SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated FROM rockets WHERE channel = ?").
					WithArgs(tc.channel).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			rocket, err := repo.GetByChannel(context.Background(), tc.channel)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Nil(t, rocket)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRocket, rocket)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRocketRepository_Save(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	now := time.Now()
	rocket := &domain.Rocket{
		Channel:     "channel-1",
		Type:        "Falcon-9",
		Speed:       1000,
		Mission:     "ARTEMIS",
		LaunchTime:  now,
		Status:      domain.RocketStatusLaunched,
		LastMessage: 1,
	}

	testCases := []struct {
		name          string
		rocket        *domain.Rocket
		expectedError string
	}{
		{
			name:          "successful_save",
			rocket:        rocket,
			expectedError: "",
		},
		{
			name:          "database_error",
			rocket:        rocket,
			expectedError: "failed to save rocket: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectExec("INSERT INTO rockets").
					WithArgs(
						tc.rocket.Channel,
						tc.rocket.Type,
						tc.rocket.Speed,
						tc.rocket.Mission,
						tc.rocket.LaunchTime,
						tc.rocket.Status,
						nil,
						tc.rocket.Reason,
						sqlmock.AnyArg(), // last_updated
						tc.rocket.LastMessage,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO rockets").
					WithArgs(
						tc.rocket.Channel,
						tc.rocket.Type,
						tc.rocket.Speed,
						tc.rocket.Mission,
						tc.rocket.LaunchTime,
						tc.rocket.Status,
						nil,
						tc.rocket.Reason,
						sqlmock.AnyArg(), // last_updated
						tc.rocket.LastMessage,
					).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			err := repo.Save(context.Background(), tc.rocket)

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

func TestRocketRepository_Update(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	now := time.Now()
	rocket := &domain.Rocket{
		Channel:     "channel-1",
		Type:        "Falcon-9",
		Speed:       1000,
		Mission:     "ARTEMIS",
		LaunchTime:  now,
		Status:      domain.RocketStatusLaunched,
		LastMessage: 1,
	}

	testCases := []struct {
		name          string
		rocket        *domain.Rocket
		expectedError string
	}{
		{
			name:          "successful_update",
			rocket:        rocket,
			expectedError: "",
		},
		{
			name:          "database_error",
			rocket:        rocket,
			expectedError: "failed to update rocket: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectExec("UPDATE rockets").
					WithArgs(
						tc.rocket.Type,
						tc.rocket.Speed,
						tc.rocket.Mission,
						tc.rocket.Status,
						nil,
						tc.rocket.Reason,
						sqlmock.AnyArg(), // last_updated
						tc.rocket.LastMessage,
						tc.rocket.Channel,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("UPDATE rockets").
					WithArgs(
						tc.rocket.Type,
						tc.rocket.Speed,
						tc.rocket.Mission,
						tc.rocket.Status,
						nil,
						tc.rocket.Reason,
						sqlmock.AnyArg(), // last_updated
						tc.rocket.LastMessage,
						tc.rocket.Channel,
					).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			err := repo.Update(context.Background(), tc.rocket)

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

func TestRocketRepository_Delete(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	testCases := []struct {
		name          string
		channel       string
		expectedError string
	}{
		{
			name:          "successful_delete",
			channel:       "channel-1",
			expectedError: "",
		},
		{
			name:          "database_error",
			channel:       "channel-1",
			expectedError: "failed to delete rocket: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectExec("DELETE FROM rockets").
					WithArgs(tc.channel).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("DELETE FROM rockets").
					WithArgs(tc.channel).
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			err := repo.Delete(context.Background(), tc.channel)

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

func TestRocketRepository_BeginTx(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	testCases := []struct {
		name          string
		expectedError string
	}{
		{
			name:          "successful_begin",
			expectedError: "",
		},
		{
			name:          "database_error",
			expectedError: "failed to begin transaction: sql: connection is already closed",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" {
				mock.ExpectBegin()
			} else {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			tx, err := repo.BeginTx(context.Background())

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRocketRepository_GetAll(t *testing.T) {
	// Create sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewRocketRepository(db)

	now := time.Now()
	explodedAt := now.Add(time.Hour)

	testCases := []struct {
		name          string
		sortBy        string
		order         string
		mockRows      *sqlmock.Rows
		expectedError string
		expectedCount int
	}{
		{
			name:   "default_sorting",
			sortBy: "",
			order:  "",
			mockRows: sqlmock.NewRows([]string{
				"channel", "type", "speed", "mission", "launch_time", "status",
				"exploded_at", "reason", "last_updated",
			}).AddRow(
				"channel-1", "type-1", 100, "mission-1", now, "launched",
				explodedAt, "reason-1", now,
			).AddRow(
				"channel-2", "type-2", 200, "mission-2", now.Add(time.Hour), "exploded",
				nil, "", now,
			),
			expectedError: "",
			expectedCount: 2,
		},
		{
			name:   "custom_sorting",
			sortBy: "speed",
			order:  "ASC",
			mockRows: sqlmock.NewRows([]string{
				"channel", "type", "speed", "mission", "launch_time", "status",
				"exploded_at", "reason", "last_updated",
			}).AddRow(
				"channel-1", "type-1", 100, "mission-1", now, "launched",
				nil, "", now,
			).AddRow(
				"channel-2", "type-2", 200, "mission-2", now, "launched",
				nil, "", now,
			),
			expectedError: "",
			expectedCount: 2,
		},
		{
			name:          "invalid_sort_column",
			sortBy:        "invalid_column",
			order:         "ASC",
			mockRows:      nil,
			expectedError: "invalid sort column: invalid_column",
			expectedCount: 0,
		},
		{
			name:          "invalid_sort_order",
			sortBy:        "speed",
			order:         "INVALID",
			mockRows:      nil,
			expectedError: "invalid sort order: INVALID",
			expectedCount: 0,
		},
		{
			name:   "database_error",
			sortBy: "",
			order:  "",
			mockRows: sqlmock.NewRows([]string{
				"channel", "type", "speed", "mission", "launch_time", "status",
				"exploded_at", "reason", "last_updated",
			}),
			expectedError: "failed to get rockets: sql: connection is already closed",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Set up expectations
			if tc.expectedError == "" && tc.mockRows != nil {
				expectedQuery := `SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated 
								FROM rockets 
								ORDER BY `
				if tc.sortBy != "" {
					expectedQuery += tc.sortBy + " " + tc.order
				} else {
					expectedQuery += "type DESC"
				}

				mock.ExpectQuery(expectedQuery).
					WillReturnRows(tc.mockRows)
			} else if tc.expectedError != "" && tc.mockRows != nil {
				mock.ExpectQuery("SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated FROM rockets").
					WillReturnError(sql.ErrConnDone)
			}

			// Execute test
			rockets, err := repo.GetAll(context.Background(), tc.sortBy, tc.order)

			// Check results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Nil(t, rockets)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rockets)
				assert.Equal(t, tc.expectedCount, len(rockets))

				// Verify rocket data for successful cases
				if tc.expectedCount > 0 {
					// Check first rocket
					assert.Equal(t, "channel-1", rockets[0].Channel)
					assert.Equal(t, "type-1", rockets[0].Type)
					assert.Equal(t, 100, rockets[0].Speed)
					assert.Equal(t, "mission-1", rockets[0].Mission)
					assert.Equal(t, "launched", rockets[0].Status)

					// For default sorting test case, check second rocket
					if tc.name == "default_sorting" {
						assert.Equal(t, "channel-2", rockets[1].Channel)
						assert.Equal(t, "type-2", rockets[1].Type)
						assert.Equal(t, 200, rockets[1].Speed)
						assert.Equal(t, "mission-2", rockets[1].Mission)
						assert.Equal(t, "exploded", rockets[1].Status)
						assert.Nil(t, rockets[1].ExplodedAt)
						assert.Empty(t, rockets[1].Reason)
					}
				}
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
