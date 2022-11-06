package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	LabelPackage string = "package"

	FormatJSON string = "json"
	FormatText string = "text"

	KindStdout string = "stdout"
)

func NewLogger(pkgName string) *log.Entry {
	if pkgName == "" {
		pkgName = "main"
	}

	logger := log.New()

	switch Format() {
	case FormatJSON:
		logger.Formatter = &log.JSONFormatter{}
	case FormatText:
		fallthrough
	default:
		logger.Formatter = &log.TextFormatter{}
	}

	if Kind() == KindStdout {
		logger.Out = os.Stdout
	}

	logger.Level = Level()

	return logger.WithField(LabelPackage, pkgName)
}
