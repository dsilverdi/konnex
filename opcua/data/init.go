package data

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := "postgres://opcua:konnexopcua@opcua-db:5432/opcua"

	var err error
	var conn *pgxpool.Pool

	for {
		conn, err = pgxpool.Connect(ctx, connStr)
		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "connect: connection refused") {
			return nil, err
		}

		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}

	//run a simple query to check our connection
	if err = execTable(ctx, conn); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func execTable(ctx context.Context, conn *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS node (
			id VARCHAR(255) PRIMARY KEY,
			channel_id VARCHAR(255),
			server_uri VARCHAR(255),
			node_id VARCHAR(255)
		);`,
		`CREATE TABLE IF NOT EXISTS node_data (
			time TIMESTAMPTZ NOT NULL,
			thing_id VARCHAR(255),
			data VARCHAR(255),
			data_type VARCHAR(255),
			FOREIGN KEY (thing_id) REFERENCES node (id) ON DELETE CASCADE
			);
			SELECT create_hypertable('node_data', 'time', if_not_exists => TRUE)`,
	}

	var err error
	for i := range queries {
		_, err = conn.Exec(ctx, queries[i])
		if err != nil {
			return err
		}
	}

	return nil
}
