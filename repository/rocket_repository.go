package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"lunar-rockets/domain"
)

type sqliteTransaction struct {
	tx *sql.Tx
}

func (t *sqliteTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *sqliteTransaction) Rollback() error {
	return t.tx.Rollback()
}

type RocketRepository struct {
	db *sql.DB
}

func NewRocketRepository(db *sql.DB) domain.RocketRepository {
	return &RocketRepository{db: db}
}

func (r *RocketRepository) GetByChannel(ctx context.Context, channel string) (*domain.Rocket, error) {
	query := `SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated 
			  FROM rockets 
			  WHERE channel = ?`

	var rocket domain.Rocket
	var explodedAt sql.NullTime
	var reason sql.NullString

	err := r.db.QueryRowContext(ctx, query, channel).Scan(
		&rocket.Channel,
		&rocket.Type,
		&rocket.Speed,
		&rocket.Mission,
		&rocket.LaunchTime,
		&rocket.Status,
		&explodedAt,
		&reason,
		&rocket.LastUpdated,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rocket: %w", err)
	}

	if explodedAt.Valid {
		t := explodedAt.Time
		rocket.ExplodedAt = &t
	}

	if reason.Valid {
		rocket.Reason = reason.String
	}

	return &rocket, nil
}

func (r *RocketRepository) GetAll(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error) {
	if sortBy == "" {
		sortBy = "type"
	}
	if order == "" {
		order = "DESC"
	}

	validColumns := map[string]bool{
		"channel": true, "type": true, "speed": true, "mission": true, "status": true,
	}

	if !validColumns[sortBy] {
		return nil, fmt.Errorf("invalid sort column: %s", sortBy)
	}

	if order != "ASC" && order != "DESC" {
		return nil, fmt.Errorf("invalid sort order: %s", order)
	}

	query := fmt.Sprintf(`SELECT channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated 
						  FROM rockets 
						  ORDER BY %s %s`, sortBy, order)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rockets: %w", err)
	}
	defer rows.Close()

	var rockets []*domain.Rocket

	for rows.Next() {
		var rocket domain.Rocket
		var explodedAt sql.NullTime
		var reason sql.NullString

		err := rows.Scan(
			&rocket.Channel,
			&rocket.Type,
			&rocket.Speed,
			&rocket.Mission,
			&rocket.LaunchTime,
			&rocket.Status,
			&explodedAt,
			&reason,
			&rocket.LastUpdated,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan rocket: %w", err)
		}

		if explodedAt.Valid {
			t := explodedAt.Time
			rocket.ExplodedAt = &t
		}

		if reason.Valid {
			rocket.Reason = reason.String
		}

		rockets = append(rockets, &rocket)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rockets: %w", err)
	}

	return rockets, nil
}

func (r *RocketRepository) Save(ctx context.Context, rocket *domain.Rocket) error {
	query := `INSERT INTO rockets (
				channel, type, speed, mission, launch_time, status, exploded_at, reason, last_updated, last_message
			  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var explodedAt interface{}
	if rocket.ExplodedAt != nil {
		explodedAt = *rocket.ExplodedAt
	}

	_, err := r.db.ExecContext(ctx, query,
		rocket.Channel,
		rocket.Type,
		rocket.Speed,
		rocket.Mission,
		rocket.LaunchTime,
		rocket.Status,
		explodedAt,
		rocket.Reason,
		time.Now(),
		rocket.LastMessage,
	)

	if err != nil {
		return fmt.Errorf("failed to save rocket: %w", err)
	}

	return nil
}

func (r *RocketRepository) Update(ctx context.Context, rocket *domain.Rocket) error {
	query := `UPDATE rockets 
			  SET type = ?, speed = ?, mission = ?, status = ?, 
				  exploded_at = ?, reason = ?, last_updated = ?, last_message = ?
			  WHERE channel = ?`

	var explodedAt interface{}
	if rocket.ExplodedAt != nil {
		explodedAt = *rocket.ExplodedAt
	}

	_, err := r.db.ExecContext(ctx, query,
		rocket.Type,
		rocket.Speed,
		rocket.Mission,
		rocket.Status,
		explodedAt,
		rocket.Reason,
		time.Now(),
		rocket.LastMessage,
		rocket.Channel,
	)

	if err != nil {
		return fmt.Errorf("failed to update rocket: %w", err)
	}

	return nil
}

func (r *RocketRepository) Delete(ctx context.Context, channel string) error {
	query := `DELETE FROM rockets WHERE channel = ?`

	_, err := r.db.ExecContext(ctx, query, channel)
	if err != nil {
		return fmt.Errorf("failed to delete rocket: %w", err)
	}

	return nil
}

func (r *RocketRepository) BeginTx(ctx context.Context) (domain.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &sqliteTransaction{tx: tx}, nil
}
