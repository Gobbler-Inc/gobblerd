package config

import (
	"fmt"
	"os"

	"github.com/alfreddobradi/go-bb-man/database/cockroach"
	"github.com/alfreddobradi/goconf"
)

var Cfg *goconf.Configuration

func Load(path string) error {
	grammar := struct {
		Database struct {
			Kind string `env:"GOBBLER_DB_KIND"`
			CRDB struct {
				// postgresql://gobb:TTjYniGOFXQK6CrK2emurA@free-tier5.gcp-europe-west1.cockroachlabs.cloud:26257/gobb-dev?application_name=ccloud&options=--cluster%3Dpurple-moose-1962&sslmode=verify-full&sslrootcert=%2FUsers%2Falfred.dobradi%2F.postgresql%2Froot.crt
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