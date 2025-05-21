package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"lunar-rockets/domain"
)

type RocketMessageUsecase interface {
	ProcessMessage(ctx context.Context, message *domain.RocketMessage) error
}

type rocketMessageUsecase struct {
	rocketRepo         domain.RocketRepository
	messageRepo        domain.MessageRepository
	rocketStateUsecase RocketStateUsecase
	messageBuffer      map[string]map[int64]*domain.RocketMessage
	bufferMutex        sync.RWMutex
}

func NewRocketMessageUsecase(rocketRepo domain.RocketRepository, messageRepo domain.MessageRepository, rocketStateUsecase RocketStateUsecase) RocketMessageUsecase {
	return &rocketMessageUsecase{
		rocketRepo:         rocketRepo,
		messageRepo:        messageRepo,
		rocketStateUsecase: rocketStateUsecase,
		messageBuffer:      make(map[string]map[int64]*domain.RocketMessage),
	}
}

func (p *rocketMessageUsecase) ProcessMessage(ctx context.Context, message *domain.RocketMessage) error {
	lastMessageNumber, err := p.messageRepo.FindLastMessageNumber(ctx, message.Metadata.Channel)
	if err != nil {
		return fmt.Errorf("failed to check if message was processed: %w", err)
	}

	// Skip processed messages
	if lastMessageNumber >= message.Metadata.MessageNumber {
		log.Printf("Skipping already processed message %d for channel %s", message.Metadata.MessageNumber, message.Metadata.Channel)
		return nil
	}

	// Buffer out-of-order messages
	if lastMessageNumber+1 < message.Metadata.MessageNumber {
		log.Printf("Buffering out-of-order message %d for channel %s", message.Metadata.MessageNumber, message.Metadata.Channel)
		p.addToBuffer(message)
		return nil
	}

	// Process message, it's the expected one
	err = p.rocketStateUsecase.UpdateRocketFromMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to execute rocket state usecase: %w", err)
	}

	if err := p.processBufferedMessages(ctx, message.Metadata.Channel, message.Metadata.MessageNumber); err != nil {
		return err
	}

	log.Printf("Successfully processed message %d for channel %s", message.Metadata.MessageNumber, message.Metadata.Channel)
	return nil
}

// addToBuffer adds a message to the buffer for its channel
func (p *rocketMessageUsecase) addToBuffer(message *domain.RocketMessage) {
	p.bufferMutex.Lock()
	defer p.bufferMutex.Unlock()

	channel := message.Metadata.Channel
	if _, exists := p.messageBuffer[channel]; !exists {
		p.messageBuffer[channel] = make(map[int64]*domain.RocketMessage)
	}
	p.messageBuffer[channel][message.Metadata.MessageNumber] = message
}

// processBufferedMessages processes consecutive messages from the buffer
func (p *rocketMessageUsecase) processBufferedMessages(ctx context.Context, channel string, lastProcessedNumber int64) error {
	p.bufferMutex.Lock()
	defer p.bufferMutex.Unlock()

	channelBuffer, exists := p.messageBuffer[channel]
	if !exists {
		return nil
	}

	nextNumber := lastProcessedNumber + 1
	for {
		message, exists := channelBuffer[nextNumber]
		if !exists {
			break
		}

		err := p.rocketStateUsecase.UpdateRocketFromMessage(ctx, message)
		if err != nil {
			return fmt.Errorf("failed to process buffered message %d: %w", nextNumber, err)
		}

		delete(channelBuffer, nextNumber)
		nextNumber++
	}

	if len(channelBuffer) == 0 {
		delete(p.messageBuffer, channel)
	}

	return nil
}
