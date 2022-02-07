package things

import (
	"context"
	"fmt"
	"konnex"
	"konnex/pkg/errors"
)

type Service interface {
	// Create Things
	CreateThings(ctx context.Context, t Things) (*Things, error)

	// Get List of Things
	GetThings(ctx context.Context) ([]Things, error)

	// Get Specific Thing
	GetSpecificThing(ctx context.Context, id string) (*Things, error)

	// Delte Thing by ID
	DeleteThing(ctx context.Context, id string) error

	// Create IoT Channel
	CreateChannel(ctx context.Context, ch Channel) (*Channel, error)

	//Get List of IoT Channel
	GetChannels(ctx context.Context) ([]Channel, error)

	//Get Specific IoT Channel
	GetSpecificChannel(ctx context.Context, id string) (*Channel, error)

	// Delete Channel by ID
	DeleteChannel(ctx context.Context, id string) error

	// Connect Channel and IoT
	// Connect(ctx context.Context, ThingsID string, ChannelID string) error

	// Disconnect Channel and IoT
	// Disconnect(ctx context.Context, ThingsID string, ChannelID string) error
}

type thingsService struct {
	ThingRepository   ThingRepository
	ChannelRepository ChannelRepository
	IDprovider        konnex.IDprovider
}

func New(trepo ThingRepository, chrepo ChannelRepository, uid konnex.IDprovider) Service {
	return &thingsService{
		ThingRepository:   trepo,
		ChannelRepository: chrepo,
		IDprovider:        uid,
	}
}

func (s *thingsService) CreateThings(ctx context.Context, t Things) (*Things, error) {
	if t.ID == "" {
		id, err := s.IDprovider.ID()
		if err != nil {
			return nil, errors.Wrap(errors.ErrCreateUUID, err)
		}

		t.ID = id
	}

	t.Owner = "annon"
	err := s.ThingRepository.Insert(ctx, t)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCreateEntity, err)
	}

	return &t, nil
}

func (s *thingsService) GetThings(ctx context.Context) ([]Things, error) {
	var thingsList []Things

	thingsList, err := s.ThingRepository.GetAll(ctx)
	if err != nil {
		fmt.Println("db error")
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return thingsList, nil
}

func (s *thingsService) GetSpecificThing(ctx context.Context, id string) (*Things, error) {
	things, err := s.ThingRepository.GetSpecific(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return things, nil
}

func (s *thingsService) DeleteThing(ctx context.Context, id string) error {
	err := s.ThingRepository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(errors.ErrRemoveEntity, err)
	}

	return nil
}

// Channel Services
func (s *thingsService) CreateChannel(ctx context.Context, ch Channel) (*Channel, error) {
	if ch.ID == "" {
		id, err := s.IDprovider.ID()
		if err != nil {
			return nil, errors.Wrap(errors.ErrCreateUUID, err)
		}

		ch.ID = id
	}

	ch.Owner = "annon"
	err := s.ChannelRepository.Insert(ctx, ch)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCreateEntity, err)
	}

	return &ch, nil
}

func (s *thingsService) GetChannels(ctx context.Context) ([]Channel, error) {
	var channels []Channel

	channels, err := s.ChannelRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return channels, nil
}

func (s *thingsService) GetSpecificChannel(ctx context.Context, id string) (*Channel, error) {
	channel, err := s.ChannelRepository.GetSpecific(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return channel, nil
}

func (s *thingsService) DeleteChannel(ctx context.Context, id string) error {
	err := s.ChannelRepository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(errors.ErrRemoveEntity, err)
	}

	return nil
}

// Connection Services

// func (s *thingsService) Connect(ctx context.Context, ThingsID string, ChannelID string) error {
// 	return nil
// }

// func (s *thingsService) Disconnect(ctx context.Context, ThingsID string, ChannelID string) error {
// 	return nil
// }
