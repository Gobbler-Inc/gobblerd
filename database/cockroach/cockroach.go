package cockroach

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alfreddobradi/go-bb-man/parser"
	"github.com/google/uuid"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	pgx "github.com/jackc/pgx/v4"
)

type DB struct {
	*pgx.Conn
}

func New() (*DB, error) {
	log.Printf("Connecting to CockroachDB at %s", Host())
	connUrl := createConnUrl()

	config, err := pgx.ParseConfig(connUrl)
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

func (db *DB) GetReplayList() ([]parser.Record, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM replays")
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve rows: %w", err)
	}
	response := make([]parser.Record, 0)
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var home string
		var away string
		if err := rows.Scan(&id, &home, &away); err != nil {
			return nil, fmt.Errorf("Failed to scan row into struct: %w", err)
		}

		var homeStruct parser.TeamStats
		var awayStruct parser.TeamStats

		if err := json.Unmarshal([]byte(home), &homeStruct); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal home team data in replay %s: %w", id.String(), err)
		}

		if err := json.Unmarshal([]byte(away), &awayStruct); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal away team data in replay %s: %w", id.String(), err)
		}

		response = append(response, parser.Record{
			ID:   id,
			Home: homeStruct,
			Away: awayStruct,
		})
	}

	return response, nil
}

func (db *DB) GetReplay(id uuid.UUID) (parser.Record, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM replays WHERE id = $1", id)
	if err != nil {
		return parser.Record{}, fmt.Errorf("Failed to retrieve rows: %w", err)
	}
	response := make([]parser.Record, 0)
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var home string
		var away string
		if err := rows.Scan(&id, &home, &away); err != nil {
			return parser.Record{}, fmt.Errorf("Failed to scan row into struct: %w", err)
		}

		var homeStruct parser.TeamStats
		var awayStruct parser.TeamStats

		if err := json.Unmarshal([]byte(home), &homeStruct); err != nil {
			return parser.Record{}, fmt.Errorf("Failed to unmarshal home team data in replay %s: %w", id.String(), err)
		}

		if err := json.Unmarshal([]byte(away), &awayStruct); err != nil {
			return parser.Record{}, fmt.Errorf("Failed to unmarshal away team data in replay %s: %w", id.String(), err)
		}

		response = append(response, parser.Record{
			ID:   id,
			Home: homeStruct,
			Away: awayStruct,
		})
	}

	return response[0], nil
}

func createConnUrl() string {
	auth := ""
	if Username() != "" {
		auth = fmt.Sprint(Username())
		if Password() != "" {
			auth = fmt.Sprintf("%s:%s", auth, Password())
		}
		auth = fmt.Sprintf("%s@", auth)
	}

	url := fmt.Sprintf("postgres://%s%s:%d/%s", auth, Host(), Port(), Database())

	params := make([]string, 0)
	if Options() != "" {
		params = append(params, fmt.Sprintf("options=%s", Options()))
	}

	if SSLMode() != "" {
		params = append(params, fmt.Sprintf("sslmode=%s", SSLMode()))
	}

	if SSLRootCert() != "" {
		params = append(params, fmt.Sprintf("sslrootcert=%s", SSLRootCert()))
	}

	if len(params) > 0 {
		url = fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
	}

	return url
}
