package data

import (
	"context"
	"konnex/opcua"

	"github.com/jackc/pgx/v4/pgxpool"
)

type NodeRepository struct {
	conn *pgxpool.Pool
}

func NewNodeRepo(db *pgxpool.Pool) opcua.NodeRepository {
	return &NodeRepository{
		conn: db,
	}
}

func (db *NodeRepository) Save(ctx context.Context, node *opcua.Node) error {
	query := `INSERT INTO node (id, channel_id, server_uri, node_id) VALUES ($1, $2, $3, $4);`

	_, err := db.conn.Exec(ctx, query, node.ID, node.ChannelID, node.ServerUri, node.NodeID)
	if err != nil {
		return err
	}

	return nil
}

func (db *NodeRepository) ReadAll(ctx context.Context) ([]opcua.Node, error) {
	var nodes []opcua.Node

	query := `SELECT id, server_uri, node_id FROM node`

	rows, err := db.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r opcua.Node
		err = rows.Scan(&r.ID, &r.ServerUri, &r.NodeID)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, r)
	}

	return nodes, nil
}

func (db *NodeRepository) ReadbyID(ctx context.Context, id string) (*opcua.Node, error) {
	var node *opcua.Node
	return node, nil
}

func (db *NodeRepository) Delete(ctx context.Context, id, tag string) error {
	var query string
	switch tag {
	case "thing":
		query = `DELETE FROM node WHERE id = $1`
	case "channel":
		query = `DELETE FROM node WHERE channel_id = $1`
	}

	_, err := db.conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
