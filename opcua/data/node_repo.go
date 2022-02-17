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
	query := `INSERT INTO node (id, server_uri, node_id) VALUES ($1, $2, $3);`

	_, err := db.conn.Exec(ctx, query, node.ID, node.ServerUri, node.NodeID)
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

func (db *NodeRepository) Delete(ctx context.Context, id string) error {
	return nil
}
