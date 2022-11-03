package database

import (
	"github.com/alfreddobradi/go-bb-man/parser"
	"github.com/google/uuid"
)

type DB interface {
	SaveReplay(record parser.Record) error
	GetReplayList() ([]parser.Record, error)
	GetReplay(id uuid.UUID) (parser.Record, error)
}
