package config

import (
	"fmt"
	"os"

	"github.com/alfreddobradi/go-bb-man/database/cockroach"
	"github.com/alfreddobradi/go-bb-man/logging"
	"github.com/alfreddobradi/goconf"
)

var Cfg *goconf.Configuration

func Load(path string) error {
	grammar := struct {
		Logging struct {
			Format string `env:"GOBBLER_LOGGING_FORMAT"`
			Kind   string `env:"GOBBLER_LOGGING_KIND"`
			Path   string `env:"GOBBLER_LOGGING_PATH"`
			Level  string `env:"GOBBLER_LOGGING_LEVEL"`
		}
		Database struct {
			Kind string `env:"GOBBLER_DB_KIND"`
			CRDB struct {
				Username    string `env:"GOBBLER_DB_USERNAME"`
				Password    string `env:"GOBBLER_DB_PASSWORD"`
				Host        string `env:"GOBBLER_DB_HOST"`
				Port        int    `env:"GOBBLER_DB_PORT"`
				Database    string `env:"GOBBLER_DB_DATABASE"`
				Options     string `env:"GOBBLER_DB_OPTIONS"`
				SSLMode     string `yaml:"ssl_mode" env:"GOBBLER_DB_SSL_MODE"`
				SSLRootCert string `yaml:"ssl_root_cert" env:"GOBBLER_DB_SSL_ROOT_CERT"`
			} `yaml:"crdb"`
		}
	}{}

	fp, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return fmt.Errorf("Failed to open config file %s: %w", path, err)
	}
	defer fp.Close()

	config, err := goconf.Load(&grammar, fp)
	if err != nil {
		return fmt.Errorf("Failed to load configuration: %w", err)
	}

	Cfg = config

	SetLoggingConfig(config)

	if config.GetString("database.kind") == "crdb" {
		SetCockroachConfig(config)
	}

	return nil
}

func SetCockroachConfig(config *goconf.Configuration) {
	if host := config.GetString("database.crdb.host"); host != "" && host != cockroach.Host() {
		cockroach.SetHost(host)
	}

	if port := config.GetInt("database.crdb.port"); port != 0 && port != cockroach.Port() {
		cockroach.SetPort(port)
	}

	if username := config.GetString("database.crdb.username"); username != "" && username != cockroach.Username() {
		cockroach.SetUsername(username)
	}

	if password := config.GetString("database.crdb.password"); password != "" && password != cockroach.Password() {
		cockroach.SetPassword(password)
	}

	if database := config.GetString("database.crdb.database"); database != "" && database != cockroach.Database() {
		cockroach.SetDatabase(database)
	}

	if options := config.GetString("database.crdb.options"); options != "" && options != cockroach.Options() {
		cockroach.SetOptions(options)
	}

	if sslMode := config.GetString("database.crdb.ssl_mode"); sslMode != "" && sslMode != cockroach.SSLMode() {
		cockroach.SetSSLMode(sslMode)
	}

	if sslRootCert := config.GetString("database.crdb.ssl_root_cert"); sslRootCert != "" && sslRootCert != cockroach.SSLRootCert() {
		cockroach.SetSSLRootCert(sslRootCert)
	}
}

func SetLoggingConfig(config *goconf.Configuration) {
	if format := config.GetString("logging.format"); format != logging.Format() {
		logging.SetFormat(format)
	}

	if kind := config.GetString("logging.kind"); kind != logging.Kind() {
		logging.SetKind(kind)
	}

	if path := config.GetString("logging.path"); path != logging.Path() {
		logging.SetPath(path)
	}

	if level := config.GetString("logging.level"); level != logging.Level().String() {
		logging.SetLevel(level)
	}
}
