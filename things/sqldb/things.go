package sqldb

import (
	"context"
	"konnex/things"
)

type ThingRepository struct {
	db Database
}

func NewThingRepository(db Database) things.ThingRepository {
	return &ThingRepository{
		db: db,
	}
}

func (t ThingRepository) Insert(ctx context.Context, things things.Things) error {
	return nil
}
