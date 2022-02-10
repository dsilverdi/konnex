package sqldb

import (
	"context"
	"encoding/json"
	"konnex/things"
)

type ChannelRepository struct {
	db Database
}
type ChannelDB struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Type     string `db:"type"`
	Metadata string `db:"metadata"`
}

func NewChannelRepository(db Database) things.ChannelRepository {
	return &ChannelRepository{
		db: db,
	}
}

func (ch ChannelRepository) Insert(ctx context.Context, channel things.Channel) error {
	query := `INSERT INTO channels (id, owner, name, type, metadata)
	VALUES (:id, :owner, :name, :type, :metadata);`

	ChDB, err := toDBChannel(channel)
	if err != nil {
		return err
	}

	_, err = ch.db.NamedExecContext(ctx, query, ChDB)
	if err != nil {
		return err
	}

	return nil
}

func (ch ChannelRepository) GetAll(ctx context.Context, owner string) ([]things.Channel, error) {
	var channelList []things.Channel

	query := `SELECT id, owner, name, type, metadata FROM channels WHERE owner = ?;`

	rows, err := ch.db.QueryxContext(ctx, query, owner)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var chDB ChannelDB
		err = rows.StructScan(&chDB)
		if err != nil {
			return nil, err
		}

		channels, err := toChannel(chDB)
		if err != nil {
			return nil, err
		}

		channelList = append(channelList, channels)
	}

	return channelList, nil
}

func (ch ChannelRepository) GetSpecific(ctx context.Context, ID, owner string) (*things.Channel, error) {
	var Channel things.Channel
	var channelDB ChannelDB

	query := `SELECT id, owner, name, type, metadata FROM channels WHERE id = ? AND owner = ?`

	err := ch.db.QueryRowxContext(ctx, query, ID, owner).StructScan(&channelDB)
	if err != nil {
		return nil, err
	}

	Channel, err = toChannel(channelDB)
	if err != nil {
		return nil, err
	}

	return &Channel, nil
}

func (ch ChannelRepository) Delete(ctx context.Context, id, owner string) error {
	dbCh := ChannelDB{
		ID:    id,
		Owner: owner,
	}

	query := `DELETE FROM channels WHERE id = :id AND owner = :owner`

	_, err := ch.db.NamedExecContext(ctx, query, dbCh)
	if err != nil {
		return err
	}

	return nil
}

func toDBChannel(ch things.Channel) (ChannelDB, error) {
	var data string
	if len(ch.Metadata) > 0 {
		b, err := json.Marshal(ch.Metadata)
		if err != nil {
			return ChannelDB{}, err
		}
		data = string(b)
	}

	return ChannelDB{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Type:     ch.Type,
		Metadata: data,
	}, nil
}

func toChannel(dbch ChannelDB) (things.Channel, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbch.Metadata), &metadata); err != nil {
		return things.Channel{}, err
	}

	return things.Channel{
		ID:       dbch.ID,
		Owner:    dbch.Owner,
		Name:     dbch.Name,
		Type:     dbch.Type,
		Metadata: metadata,
	}, nil
}
