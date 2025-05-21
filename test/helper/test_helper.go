package helper

import (
	"time"

	"lunar-rockets/domain"
)

// CreateTestRocket creates a rocket with given parameters for testing purposes
func CreateTestRocket(channel, rocketType, mission, status string, speed int, launchTime time.Time) *domain.Rocket {
	return &domain.Rocket{
		Channel:     channel,
		Type:        rocketType,
		Speed:       speed,
		Mission:     mission,
		LaunchTime:  launchTime,
		Status:      status,
		LastUpdated: time.Now(),
		LastMessage: 1,
	}
}

// CreateExplodedRocket creates a rocket that has exploded
func CreateExplodedRocket(channel, rocketType, mission, reason string, explodedAt time.Time) *domain.Rocket {
	rocket := CreateTestRocket(channel, rocketType, mission, "EXPLODED", 0, explodedAt.Add(-1*time.Hour))
	rocket.Reason = reason
	rocket.ExplodedAt = &explodedAt
	return rocket
}

// CreateSampleRocketList creates a list of sample rockets for testing
func CreateSampleRocketList() []*domain.Rocket {
	now := time.Now()

	return []*domain.Rocket{
		CreateTestRocket("channel-1", "Falcon-9", "ARTEMIS", "ACTIVE", 1000, now.Add(-1*time.Hour)),
		CreateTestRocket("channel-2", "Starship", "MARS", "ACTIVE", 2000, now.Add(-2*time.Hour)),
		CreateExplodedRocket("channel-3", "Saturn-V", "APOLLO", "PRESSURE_FAILURE", now.Add(-30*time.Minute)),
		CreateTestRocket("channel-4", "Falcon-Heavy", "STARLINK", "ACTIVE", 1500, now),
	}
}

// CreateTestMessage creates a message with given parameters for testing purposes
func CreateTestMessage(channel string, messageType string, messageNumber int64, timestamp time.Time) *domain.RocketMessage {
	return &domain.RocketMessage{
		Metadata: domain.MessageMetadata{
			Channel:       channel,
			MessageType:   messageType,
			MessageNumber: messageNumber,
			MessageTime:   timestamp,
		},
		Message: domain.RocketLaunchedMessage{
			Type:        "Falcon-9",
			LaunchSpeed: 1000,
			Mission:     "ARTEMIS",
		},
	}
}
