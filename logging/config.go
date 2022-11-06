package logging

import (
	log "github.com/sirupsen/logrus"
)

var (
	format string    = FormatText
	kind   string    = KindStdout
	path   string    = ""
	level  log.Level = log.InfoLevel
)

func Format() string {
	return format
}

func Kind() string {
	return kind
}

func Path() string {
	return path
}

func Level() log.Level {
	return level
}

func SetFormat(newFormat string) {
	switch newFormat {
	case FormatJSON, FormatText:
		format = newFormat
	}
}

func SetKind(newKind string) {
	// TODO implement file-backed logging
	switch newKind {
	case KindStdout:
		kind = newKind
	}
}

func SetPath(newPath string) {
	path = newPath
}

func SetLevel(newLevel string) {
	l, err := log.ParseLevel(newLevel)
	if err == nil {
		level = l
	}
}
