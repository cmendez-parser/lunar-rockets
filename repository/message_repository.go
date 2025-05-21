package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) MarkAsProcessed(ctx context.Context, channel string, messageNumber int64) error {
	query := `INSERT INTO processed_messages (channel, message_number, processed_at)
			  VALUES (?, ?, CURRENT_TIMESTAMP)`

	_, err := r.db.ExecContext(ctx, query, channel, messageNumber)
	if err != nil {
		return fmt.Errorf("failed to mark message as processed: %w", err)
	}

	return nil
}

func (r *MessageRepository) FindLastMessageNumber(ctx context.Context, channel string) (int64, error) {
	query := `SELECT MAX(message_number) FROM processed_messages WHERE channel = ?`

	var lastMessageNumber sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, channel).Scan(&lastMessageNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No messages found, return 0
		}
		return 0, fmt.Errorf("failed to find last message number: %w", err)
	}

	if !lastMessageNumber.Valid {
		return 0, nil
	}

	return lastMessageNumber.Int64, nil
}
