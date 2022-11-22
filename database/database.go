package database

import (
	"github.com/gobbler-inc/gobblerd/parser"
	"github.com/google/uuid"
)

type DB interface {
	SaveReplay(record parser.Record) error
	GetReplayList() ([]parser.Record, error)
	GetReplay(id uuid.UUID) (parser.Record, error)
}
