package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"lunar-rockets/domain"
)

type RocketStateUsecase interface {
	UpdateRocketFromMessage(ctx context.Context, message *domain.RocketMessage) error
}

type rocketStateUsecase struct {
	rocketRepo  domain.RocketRepository
	messageRepo domain.MessageRepository
}

func NewRocketStateUsecase(rocketRepo domain.RocketRepository, messageRepo domain.MessageRepository) RocketStateUsecase {
	return &rocketStateUsecase{
		rocketRepo:  rocketRepo,
		messageRepo: messageRepo,
	}
}

func (u *rocketStateUsecase) UpdateRocketFromMessage(ctx context.Context, message *domain.RocketMessage) error {
	tx, err := u.rocketRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			log.Printf("Rolling back transaction for channel %s due to error: %v", message.Metadata.Channel, err)
			_ = tx.Rollback()
		}
	}()

	var processErr error
	switch message.Metadata.MessageType {
	case domain.TypeRocketLaunched:
		processErr = u.handleRocketLaunched(ctx, message)
	case domain.TypeRocketSpeedIncreased:
		processErr = u.handleRocketSpeedIncreased(ctx, message)
	case domain.TypeRocketSpeedDecreased:
		processErr = u.handleRocketSpeedDecreased(ctx, message)
	case domain.TypeRocketExploded:
		processErr = u.handleRocketExploded(ctx, message)
	case domain.TypeRocketMissionChanged:
		processErr = u.handleRocketMissionChanged(ctx, message)
	default:
		processErr = fmt.Errorf("unknown message type: %s", message.Metadata.MessageType)
	}

	if processErr != nil {
		return fmt.Errorf("failed to update rocket state: %w", processErr)
	}

	if err = u.messageRepo.MarkAsProcessed(ctx, message.Metadata.Channel, message.Metadata.MessageNumber); err != nil {
		return fmt.Errorf("failed to mark message as processed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully updated rocket state for message type %s, channel %s", message.Metadata.MessageType, message.Metadata.Channel)
	return nil
}

func (p *rocketStateUsecase) handleRocketLaunched(ctx context.Context, message *domain.RocketMessage) error {
	var launchMsg domain.RocketLaunchedMessage
	if err := parseMessagePayload(message.Message, &launchMsg); err != nil {
		return err
	}

	rocket, err := p.rocketRepo.GetByChannel(ctx, message.Metadata.Channel)
	if err != nil {
		return err
	}

	if rocket != nil {
		log.Printf("Rocket already exists for channel %s, skipping launch", message.Metadata.Channel)
		return nil
	}

	newRocket := &domain.Rocket{
		Channel:     message.Metadata.Channel,
		Type:        launchMsg.Type,
		Speed:       launchMsg.LaunchSpeed,
		Mission:     launchMsg.Mission,
		LaunchTime:  message.Metadata.MessageTime,
		Status:      domain.RocketStatusLaunched,
		LastUpdated: time.Now(),
		LastMessage: message.Metadata.MessageNumber,
	}

	if err := p.rocketRepo.Save(ctx, newRocket); err != nil {
		return err
	}

	log.Printf("Successfully launched rocket for channel %s", message.Metadata.Channel)
	return nil
}

func (p *rocketStateUsecase) handleRocketSpeedIncreased(ctx context.Context, message *domain.RocketMessage) error {
	var speedMsg domain.RocketSpeedIncreasedMessage
	if err := parseMessagePayload(message.Message, &speedMsg); err != nil {
		return err
	}

	rocket, err := p.rocketRepo.GetByChannel(ctx, message.Metadata.Channel)
	if err != nil {
		return err
	}

	if rocket == nil {
		return fmt.Errorf("rocket not found: %s", message.Metadata.Channel)
	}

	if rocket.Status == domain.RocketStatusExploded {
		return nil
	}

	rocket.Speed += speedMsg.By
	if rocket.Speed < 0 {
		rocket.Speed = 0
	}
	rocket.LastUpdated = time.Now()
	rocket.LastMessage = message.Metadata.MessageNumber

	return p.rocketRepo.Update(ctx, rocket)
}

func (p *rocketStateUsecase) handleRocketSpeedDecreased(ctx context.Context, message *domain.RocketMessage) error {
	var speedMsg domain.RocketSpeedDecreasedMessage
	if err := parseMessagePayload(message.Message, &speedMsg); err != nil {
		return err
	}

	rocket, err := p.rocketRepo.GetByChannel(ctx, message.Metadata.Channel)
	if err != nil {
		return err
	}

	if rocket == nil {
		return fmt.Errorf("rocket not found: %s", message.Metadata.Channel)
	}

	if rocket.Status == domain.RocketStatusExploded {
		return nil
	}

	if speedMsg.By > rocket.Speed {
		rocket.Speed = 0
	} else {
		rocket.Speed -= speedMsg.By
	}
	rocket.LastUpdated = time.Now()
	rocket.LastMessage = message.Metadata.MessageNumber

	return p.rocketRepo.Update(ctx, rocket)
}

func (p *rocketStateUsecase) handleRocketExploded(ctx context.Context, message *domain.RocketMessage) error {
	var explodeMsg domain.RocketExplodedMessage
	if err := parseMessagePayload(message.Message, &explodeMsg); err != nil {
		return err
	}

	rocket, err := p.rocketRepo.GetByChannel(ctx, message.Metadata.Channel)
	if err != nil {
		return err
	}

	if rocket == nil {
		return fmt.Errorf("rocket not found: %s", message.Metadata.Channel)
	}

	if rocket.Status == domain.RocketStatusExploded {
		return nil
	}

	rocket.Status = domain.RocketStatusExploded
	rocket.Reason = explodeMsg.Reason
	explodeTime := message.Metadata.MessageTime
	rocket.ExplodedAt = &explodeTime
	rocket.LastUpdated = time.Now()
	rocket.LastMessage = message.Metadata.MessageNumber

	return p.rocketRepo.Update(ctx, rocket)
}

func (p *rocketStateUsecase) handleRocketMissionChanged(ctx context.Context, message *domain.RocketMessage) error {
	var missionMsg domain.RocketMissionChangedMessage
	if err := parseMessagePayload(message.Message, &missionMsg); err != nil {
		return err
	}

	rocket, err := p.rocketRepo.GetByChannel(ctx, message.Metadata.Channel)
	if err != nil {
		return err
	}

	if rocket == nil {
		return fmt.Errorf("rocket not found: %s", message.Metadata.Channel)
	}

	if rocket.Status == domain.RocketStatusExploded {
		return nil
	}

	rocket.Mission = missionMsg.NewMission
	rocket.LastUpdated = time.Now()
	rocket.LastMessage = message.Metadata.MessageNumber

	return p.rocketRepo.Update(ctx, rocket)
}

func parseMessagePayload(payload interface{}, dest interface{}) error {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	if err := json.Unmarshal(payloadBytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal message payload: %w", err)
	}

	return nil
}
