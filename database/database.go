package database

import "github.com/alfreddobradi/go-bb-man/parser"

type DB interface {
	SaveReplay(record parser.Record) error
}
