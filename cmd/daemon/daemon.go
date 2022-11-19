package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gobbler-inc/gobblerd/api"
	"github.com/gobbler-inc/gobblerd/config"
	"github.com/gobbler-inc/gobblerd/database/cockroach"
	"github.com/gobbler-inc/gobblerd/helper"
	"github.com/gobbler-inc/gobblerd/logging"
	"github.com/gobbler-inc/gobblerd/processor"
	"github.com/gobbler-inc/gobblerd/ui"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var (
	configPath string = "/etc/gobblerd/config.yml"
)

func main() {
	flag.StringVar(&configPath, "cfg", "/etc/gobblerd/config.yml", "Path to the config file")
	flag.Parse()

	if err := config.Load(configPath); err != nil {
		panic(err)
	}

	logger := logging.NewLogger("main")

	retries := 8
	retryInterval := 1000

	var db *cockroach.DB
	var err error

	try := 1
	for {
		db, err = cockroach.New()
		if err != nil && try == retries {
			break
		}
		if err == nil {
			break
		}
		try++
		wait := time.Duration(retryInterval) * time.Millisecond
		logger.WithError(err).WithFields(log.Fields{
			"retry_timeout": wait.String(),
			"try":           try,
		}).Warn("Connection failed")
		time.Sleep(wait)
		retryInterval = retryInterval << 2
	}

	if err != nil {
		logger.WithError(err).WithField("max_retries", retries).Fatalf("Maximum number of retries reached, giving up.")
	}
	defer db.Close(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	reg := processor.NewRegistry(db, wg)

	r := mux.NewRouter()

	r.HandleFunc("/upload", reg.HandleProcessRequest).Methods(http.MethodPost)
	r.HandleFunc("/upload", helper.CorsHandler).Methods(http.MethodOptions)

	r.HandleFunc("/api/replays", api.ReplayListHandler(db)).Methods(http.MethodGet)
	r.HandleFunc("/api/replays", helper.CorsHandler).Methods(http.MethodOptions)

	r.HandleFunc("/api/replays/{id}", api.ReplayHandler(db)).Methods(http.MethodGet)
	r.HandleFunc("/api/replays/{id}", helper.CorsHandler).Methods(http.MethodOptions)

	spaHandler := ui.NewSpaHandler()
	r.PathPrefix("/").Handler(spaHandler)

	s := http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("Error in HTTP listener")
		}
	}()

	<-sigChan
	logger.Debug("Received stop signal")
	s.Shutdown(context.Background())
	reg.Stop()
	wg.Wait()
}
