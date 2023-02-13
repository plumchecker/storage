package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/D3vR4pt0rs/logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/plumchecker/storage/internal/infrastructure/postgres"
	"github.com/plumchecker/storage/internal/interfaces/handlers"
	"github.com/plumchecker/storage/internal/interfaces/repository"
	"github.com/plumchecker/storage/internal/usecases/storage"
)

func main() {
	var config = &postgres.Config{
		IP:       "postgres",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Database: "leaks",
	}

	dbClient, err := postgres.New(config)
	if err != nil {
		log.Fatalln("Error during Postgres initialization")
	}

	leakStorage, err := storage.New(repository.New(dbClient))

	router := mux.NewRouter()
	handlers.Make(router, leakStorage)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	srv := &http.Server{
		Addr:    ":30001",
		Handler: c.Handler(router),
	}

	go func() {
		listener := make(chan os.Signal, 1)
		signal.Notify(listener, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Received a shutdown signal:", <-listener)

		if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to gracefully shutdown ", err)
		}
	}()

	logger.Info.Println("[*]  Listening...")
	if err := srv.ListenAndServe(); err != nil {
		logger.Error.Println("Failed to listen and serve ", err)
	}

	logger.Critical.Println("Server shutdown")
}
