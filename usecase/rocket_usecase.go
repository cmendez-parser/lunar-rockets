package usecase

import (
	"context"
	"fmt"
	"log"

	"lunar-rockets/domain"
)

type RocketUseCase interface {
	GetRocket(ctx context.Context, channel string) (*domain.Rocket, error)
	ListRockets(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error)
}

type rocketUseCase struct {
	rocketRepo domain.RocketRepository
}

func NewRocketUseCase(rocketRepo domain.RocketRepository) RocketUseCase {
	return &rocketUseCase{rocketRepo: rocketRepo}
}

func (u *rocketUseCase) GetRocket(ctx context.Context, channel string) (*domain.Rocket, error) {
	rocket, err := u.rocketRepo.GetByChannel(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to get rocket: %w", err)
	}

	if rocket == nil {
		return nil, domain.ErrRocketNotFound
	}

	log.Printf("Successfully retrieved rocket for channel %s", channel)
	return rocket, nil
}

func (u *rocketUseCase) ListRockets(ctx context.Context, sortBy string, order string) ([]*domain.Rocket, error) {
	if sortBy == "" {
		sortBy = "type"
	}

	if order == "" || (order != "ASC" && order != "DESC") {
		order = "DESC"
	}

	rockets, err := u.rocketRepo.GetAll(ctx, sortBy, order)
	if err != nil {
		return nil, fmt.Errorf("failed to list rockets: %w", err)
	}

	log.Printf("Successfully listed %d rockets", len(rockets))
	return rockets, nil
}
