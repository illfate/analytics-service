package analytics

import (
	"context"
	"fmt"
)

type Repository interface {
	InsertEvents(ctx context.Context, events ...Event) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvents(ctx context.Context, events ...Event) error {
	err := s.repo.InsertEvents(ctx, events...)
	if err != nil {
		return fmt.Errorf("failed to insert events: %w", err)
	}
	return nil
}
