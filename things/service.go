package things

import (
	"context"
	"database/sql"
	"fmt"
	"konnex"
	"konnex/pkg/errors"
)

type Service interface {
	// Create Things
	CreateThings(ctx context.Context, t Things, token string) (*Things, error)

	// Get List of Things
	GetThings(ctx context.Context, token string) ([]Things, error)

	// Get Specific Thing
	GetSpecificThing(ctx context.Context, id, token string) (*Things, error)

	// Delte Thing by ID
	DeleteThing(ctx context.Context, id, token string) error

	// Create IoT Channel
	CreateChannel(ctx context.Context, ch Channel, token string) (*Channel, error)

	//Get List of IoT Channel
	GetChannels(ctx context.Context, token string) ([]Channel, error)

	//Get Specific IoT Channel
	GetSpecificChannel(ctx context.Context, id, token string) (*Channel, error)

	// Delete Channel by ID
	DeleteChannel(ctx context.Context, id, token string) error

	// Connect Channel and IoT
	// Connect(ctx context.Context, ThingsID string, ChannelID string) error

	// Disconnect Channel and IoT
	// Disconnect(ctx context.Context, ThingsID string, ChannelID string) error
}

type thingsService struct {
	Auth              konnex.AuthServiceClient
	ThingRepository   ThingRepository
	ChannelRepository ChannelRepository
	IDprovider        konnex.IDprovider
}

func New(trepo ThingRepository, chrepo ChannelRepository, uid konnex.IDprovider, auth konnex.AuthServiceClient) Service {
	return &thingsService{
		Auth:              auth,
		ThingRepository:   trepo,
		ChannelRepository: chrepo,
		IDprovider:        uid,
	}
}

func (s *thingsService) CreateThings(ctx context.Context, t Things, token string) (*Things, error) {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	if t.ID == "" {
		id, err := s.IDprovider.ID()
		if err != nil {
			return nil, errors.Wrap(errors.ErrCreateUUID, err)
		}

		t.ID = id
	}

	t.Owner = user.Username
	err = s.ThingRepository.Insert(ctx, t)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCreateEntity, err)
	}

	return &t, nil
}

func (s *thingsService) GetThings(ctx context.Context, token string) ([]Things, error) {
	var thingsList []Things

	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	thingsList, err = s.ThingRepository.GetAll(ctx, user.Username)
	if err != nil {
		fmt.Println("db error")
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return thingsList, nil
}

func (s *thingsService) GetSpecificThing(ctx context.Context, id, token string) (*Things, error) {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	things, err := s.ThingRepository.GetSpecific(ctx, id, user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return things, nil
}

func (s *thingsService) DeleteThing(ctx context.Context, id, token string) error {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}

	err = s.ThingRepository.Delete(ctx, id, user.Username)
	if err != nil {
		return errors.Wrap(errors.ErrRemoveEntity, err)
	}

	return nil
}

// Channel Services
func (s *thingsService) CreateChannel(ctx context.Context, ch Channel, token string) (*Channel, error) {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	if ch.ID == "" {
		id, err := s.IDprovider.ID()
		if err != nil {
			return nil, errors.Wrap(errors.ErrCreateUUID, err)
		}

		ch.ID = id
	}

	ch.Owner = user.Username
	err = s.ChannelRepository.Insert(ctx, ch)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCreateEntity, err)
	}

	return &ch, nil
}

func (s *thingsService) GetChannels(ctx context.Context, token string) ([]Channel, error) {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	var channels []Channel

	channels, err = s.ChannelRepository.GetAll(ctx, user.Username)
	if err != nil {
		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return channels, nil
}

func (s *thingsService) GetSpecificChannel(ctx context.Context, id, token string) (*Channel, error) {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrAuthorization, err)
	}

	channel, err := s.ChannelRepository.GetSpecific(ctx, id, user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(errors.ErrViewEntity, err)
	}

	return channel, nil
}

func (s *thingsService) DeleteChannel(ctx context.Context, id, token string) error {
	auth, err := s.Auth.Authorize(ctx, &konnex.Token{Value: token})
	if err != nil {
		return errors.Wrap(errors.ErrUnauthorizedAccess, err)
	}

	user, err := s.Auth.Identify(ctx, &konnex.UserID{Value: auth.UserID})
	if err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}

	err = s.ChannelRepository.Delete(ctx, id, user.Username)
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
