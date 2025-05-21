package domain

import (
	"context"
	"time"
)

const (
	TypeRocketLaunched       = "RocketLaunched"
	TypeRocketSpeedIncreased = "RocketSpeedIncreased"
	TypeRocketSpeedDecreased = "RocketSpeedDecreased"
	TypeRocketExploded       = "RocketExploded"
	TypeRocketMissionChanged = "RocketMissionChanged"
)

type MessageMetadata struct {
	Channel       string    `json:"channel"`
	MessageNumber int64     `json:"messageNumber"`
	MessageTime   time.Time `json:"messageTime"`
	MessageType   string    `json:"messageType"`
}

type RocketMessage struct {
	Metadata MessageMetadata `json:"metadata"`
	Message  interface{}     `json:"message"`
}

type RocketLaunchedMessage struct {
	Type        string `json:"type"`
	LaunchSpeed int    `json:"launchSpeed"`
	Mission     string `json:"mission"`
}

type RocketSpeedIncreasedMessage struct {
	By int `json:"by"`
}

type RocketSpeedDecreasedMessage struct {
	By int `json:"by"`
}

type RocketExplodedMessage struct {
	Reason string `json:"reason"`
}

type RocketMissionChangedMessage struct {
	NewMission string `json:"newMission"`
}

type MessageRepository interface {
	MarkAsProcessed(ctx context.Context, channel string, messageNumber int64) error
	FindLastMessageNumber(ctx context.Context, channel string) (int64, error)
}
