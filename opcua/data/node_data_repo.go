package data

import (
	"context"
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
	query := `INSERT INTO node_data (time, thing_id, data, data_type) VALUES ($1, $2, $3, $4);`

	_, err := db.conn.Exec(ctx, query, data.Time, data.ThingID, data.Data, data.DataType)
	if err != nil {
		return err
	}

	return nil
}

func (db *NodeDataRepository) ReadbyID(ctx context.Context, id string) ([]opcua.NodeData, error) {
	var datas []opcua.NodeData
	// 	SELECT time_bucket('15 minutes', time) AS fifteen_min,
	//     location, COUNT(*),
	//     MAX(temperature) AS max_temp,
	//     MAX(humidity) AS max_hum
	//   FROM conditions
	//   WHERE time > NOW() - INTERVAL '3 hours'
	//   GROUP BY fifteen_min, location
	//   ORDER BY fifteen_min DESC, max_temp DESC;

	query := `SELECT time, thing_id, data, data_type FROM node_data
	WHERE time > NOW() - INTERVAL '7 days' AND thing_id = $1
	ORDER BY time DESC`

	rows, err := db.conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r opcua.NodeData
		err = rows.Scan(&r.Time, &r.ThingID, &r.Data, &r.DataType)
		if err != nil {
			return nil, err
		}
		datas = append(datas, r)
	}

	return datas, nil
}
