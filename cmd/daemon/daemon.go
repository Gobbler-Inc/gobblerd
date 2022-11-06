package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/alfreddobradi/go-bb-man/api"
	"github.com/alfreddobradi/go-bb-man/config"
	"github.com/alfreddobradi/go-bb-man/database/cockroach"
	"github.com/alfreddobradi/go-bb-man/helper"
	"github.com/alfreddobradi/go-bb-man/processor"
	"github.com/alfreddobradi/go-bb-man/ui"

	"github.com/gorilla/mux"
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
		log.Printf("Connection failed, retrying in %s - Error: %v", wait.String(), err)
		time.Sleep(wait)
		retryInterval = retryInterval << 2
	}

	if err != nil {
		panic(err)
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

	assetHandler := ui.NewAssetHandler("/")
	r.PathPrefix(assetHandler.Prefix).Handler(assetHandler)

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
			log.Printf("Error in HTTP listener: %v", err)
		}
	}()

	<-sigChan
	log.Println("Received signal")
	s.Shutdown(context.Background())
	reg.Stop()
	wg.Wait()
}
