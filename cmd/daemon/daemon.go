package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/alfreddobradi/go-bb-man/database/cockroach"
	"github.com/alfreddobradi/go-bb-man/processor"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	connUrl := os.Getenv("DATABASE_URL")
	db, err := cockroach.New(connUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	reg := processor.NewRegistry(db, wg)

	r.HandleFunc("/upload", reg.HandleProcessRequest).Methods(http.MethodPost)

	s := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
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
