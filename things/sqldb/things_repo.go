package sqldb

import (
	"context"
	"encoding/json"
	"konnex/things"
)

type ThingRepository struct {
	db Database
}

type ThingDB struct {
	ID        string `db:"id"`
	ChannelID string `db:"channel_id"`
	Owner     string `db:"owner"`
	Name      string `db:"name"`
	Metadata  string `db:"metadata"`
}

func NewThingRepository(db Database) things.ThingRepository {
	return &ThingRepository{
		db: db,
	}
}

func (t ThingRepository) Insert(ctx context.Context, things things.Things) error {
	query := `INSERT INTO things (id, owner, name, channel_id, metadata)
	VALUES (:id, :owner, :name, :channel_id, :metadata);`

	ThDB, err := toDBThing(things)
	if err != nil {
		return err
	}

	_, err = t.db.NamedExecContext(ctx, query, ThDB)
	if err != nil {
		return err
	}

	return nil
}

func (t ThingRepository) GetAll(ctx context.Context) ([]things.Things, error) {
	var thingsList []things.Things

	query := `SELECT id, channel_id, owner, name, metadata FROM things;`

	rows, err := t.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var thDB ThingDB
		err = rows.StructScan(&thDB)
		if err != nil {
			return nil, err
		}

		things, err := toThing(thDB)
		if err != nil {
			return nil, err
		}

		thingsList = append(thingsList, things)
	}

	return thingsList, nil
}

func (t ThingRepository) GetSpecific(ctx context.Context, ID string) (*things.Things, error) {
	var Thing things.Things
	var thingDB ThingDB

	query := `SELECT id, channel_id, owner, name, metadata FROM things WHERE id = ?`

	err := t.db.QueryRowxContext(ctx, query, ID).StructScan(&thingDB)
	if err != nil {
		return nil, err
	}

	Thing, err = toThing(thingDB)
	if err != nil {
		return nil, err
	}

	return &Thing, nil
}

func (t ThingRepository) Delete(ctx context.Context, id string) error {
	dbTh := ThingDB{
		ID: id,
	}

	query := `DELETE FROM things WHERE id = :id`

	_, err := t.db.NamedExecContext(ctx, query, dbTh)
	if err != nil {
		return err
	}

	return nil
}

func toDBThing(th things.Things) (ThingDB, error) {
	var data string
	if len(th.MetaData) > 0 {
		b, err := json.Marshal(th.MetaData)
		if err != nil {
			return ThingDB{}, err
		}
		data = string(b)
	}

	return ThingDB{
		ID:        th.ID,
		ChannelID: th.ChannelID,
		Owner:     th.Owner,
		Name:      th.Name,
		Metadata:  data,
	}, nil
}

func toThing(dbth ThingDB) (things.Things, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return things.Things{}, err
	}

	return things.Things{
		ID:        dbth.ID,
		ChannelID: dbth.ChannelID,
		Owner:     dbth.Owner,
		Name:      dbth.Name,
		MetaData:  metadata,
	}, nil
}
