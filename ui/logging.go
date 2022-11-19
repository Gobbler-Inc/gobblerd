package ui

import (
	"github.com/gobbler-inc/gobblerd/logging"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	logger = logging.NewLogger("ui")
}
