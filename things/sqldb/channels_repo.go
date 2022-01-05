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

func (ch ChannelRepository) GetAll(ctx context.Context) ([]things.Channel, error) {
	var channelList []things.Channel

	query := `SELECT id, owner, name, type, metadata FROM channels;`

	rows, err := ch.db.QueryxContext(ctx, query)
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
