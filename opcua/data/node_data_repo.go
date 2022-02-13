package data

import (
	"context"
	"fmt"
	"konnex/opcua"

	"github.com/jackc/pgx/v4/pgxpool"
)

type NodeDataRepository struct {
	conn *pgxpool.Pool
}

func NewNodeDataRepo(db *pgxpool.Pool) opcua.NodeDataRepository {
	return &NodeDataRepository{
		conn: db,
	}
}

func (db *NodeDataRepository) Save(ctx context.Context, data *opcua.NodeData) error {
	fmt.Print("JALAN DI DB | ", data.ThingID)

	query := `INSERT INTO node_data (time, thing_id, data, data_type) VALUES ($1, $2, $3, $4);`

	_, err := db.conn.Exec(ctx, query, data.Time, data.ThingID, data.Data, data.DataType)
	if err != nil {
		return err
	}

	return nil
}

func (db *NodeDataRepository) ReadbyID(ctx context.Context, id string) ([]opcua.NodeData, error) {
	var datas []opcua.NodeData

	return datas, nil
}
