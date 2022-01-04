package things

import (
	"context"
)

type Service interface {
	AddThings(ctx context.Context, t Things) error
	GetThings(ctx context.Context) error
}

type thingsService struct {
	ThingRepository ThingRepository
}

func New(trepo ThingRepository) Service {
	return &thingsService{
		ThingRepository: trepo,
	}
}

func (s *thingsService) AddThings(ctx context.Context, t Things) error {
	return nil
}

func (s *thingsService) GetThings(ctx context.Context) error {
	return nil
}
