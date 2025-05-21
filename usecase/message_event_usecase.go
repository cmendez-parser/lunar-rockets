package usecase

import (
	"context"
	"fmt"

	"lunar-rockets/domain"
)

type MessageEventUsecase interface {
	ProcessMessage(ctx context.Context, message *domain.RocketMessage) error
}

type messageEventUsecase struct {
	rocketRepo         domain.RocketRepository
	messageRepo        domain.MessageRepository
	rocketStateUsecase RocketStateUsecase
}

func NewMessageEventUsecase(rocketRepo domain.RocketRepository, messageRepo domain.MessageRepository, rocketStateUsecase RocketStateUsecase) MessageEventUsecase {
	return &messageEventUsecase{
		rocketRepo:         rocketRepo,
		messageRepo:        messageRepo,
		rocketStateUsecase: rocketStateUsecase,
	}
}

func (p *messageEventUsecase) ProcessMessage(ctx context.Context, message *domain.RocketMessage) error {
	lastMessageNumber, err := p.messageRepo.FindLastMessageNumber(ctx, message.Metadata.Channel)
	if err != nil {
		return fmt.Errorf("failed to check if message was processed: %w", err)
	}

	if lastMessageNumber == 0 && message.Metadata.MessageType != domain.TypeRocketLaunched {
		return fmt.Errorf("message out of order: %w", err)
	} else if lastMessageNumber >= message.Metadata.MessageNumber {
		return nil
	}

	err = p.rocketStateUsecase.UpdateRocketFromMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to execute rocket state usecase: %w", err)
	}

	return nil
}
