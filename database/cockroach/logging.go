package cockroach

import (
	"github.com/alfreddobradi/go-bb-man/logging"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	logger = logging.NewLogger("crdb")
}
