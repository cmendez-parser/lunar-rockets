package domain

import (
	"context"
	"errors"
	"time"
)

const (
	RocketStatusLaunched = "Launched"
	RocketStatusExploded = "Exploded"
)

var (
	ErrRocketNotFound = errors.New("rocket not found")
)

type Rocket struct {
	Channel     string     `json:"channel"`              // Unique identifier for the rocket
	Type        string     `json:"type"`                 // Type of rocket
	Speed       int        `json:"speed"`                // Current speed of the rocket
	Mission     string     `json:"mission"`              // Current mission
	LaunchTime  time.Time  `json:"launchTime"`           // Time when the rocket was launched
	Status      string     `json:"status"`               // Current status
	ExplodedAt  *time.Time `json:"explodedAt,omitempty"` // Time when the rocket exploded, if applicable
	Reason      string     `json:"reason,omitempty"`     // Reason for explosion, if applicable
	LastUpdated time.Time  `json:"lastUpdated"`          // Last time the rocket state was updated
	LastMessage int64      `json:"lastMessage"`          // Last message number processed
}

type RocketRepository interface {
	GetByChannel(ctx context.Context, channel string) (*Rocket, error)
	GetAll(ctx context.Context, sortBy string, order string) ([]*Rocket, error)
	Save(ctx context.Context, rocket *Rocket) error
	Update(ctx context.Context, rocket *Rocket) error
	Delete(ctx context.Context, channel string) error
	BeginTx(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}
