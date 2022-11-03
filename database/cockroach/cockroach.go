package cockroach

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alfreddobradi/go-bb-man/parser"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	pgx "github.com/jackc/pgx/v4"
)

type DB struct {
	*pgx.Conn
}

func New(conn_url string) (*DB, error) {
	config, err := pgx.ParseConfig(conn_url)
	if err != nil {
		return nil, fmt.Errorf("Error parsing connection url: %v", err)
	}
	config.RuntimeParams["application_name"] = "$ gobb"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to cluster: %v", err)
	}

	return &DB{conn}, nil
}

func (db *DB) SaveReplay(record parser.Record) error {
	homeJson, err := json.Marshal(record.Home)
	if err != nil {
		return fmt.Errorf("Failed to marshal home team data: %v", err)
	}
	awayJson, err := json.Marshal(record.Away)
	if err != nil {
		return fmt.Errorf("Failed to marshal away team data: %v", err)
	}

	txErr := crdbpgx.ExecuteTx(context.Background(), db, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), `INSERT INTO replays (id, home_team, away_team) VALUES ($1, $2, $3)`, record.ID.String(), string(homeJson), string(awayJson))
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return fmt.Errorf("Error executing statement: %v", err)
	}

	return nil
}
