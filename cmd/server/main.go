package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"quotes/api"
	"quotes/internal/config"
	"quotes/internal/drivers"
	"quotes/internal/services"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connString := cfg.DatabaseURL
	if connString == "" {
		log.Fatal("Database connection string not set")
	}

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	driver := drivers.NewQuoteDriver(dbpool)
	service := services.NewQuoteService(driver)
	controller := api.NewQuoteController(service)

	router := mux.NewRouter()
	controller.RegisterRoutes(router)

	addr := ":" + cfg.Port
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
