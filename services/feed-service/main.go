package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adrisongomez/cqrs-events-golang/database"
	"github.com/adrisongomez/cqrs-events-golang/events"
	"github.com/adrisongomez/cqrs-events-golang/repository"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf(
		"postgres://%s:%s@postgres/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
	)

	repo, err := database.NewPosgrestRepository(addr)

	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepository(repo)

	n, err := events.NewNatsEventStore(
		fmt.Sprintf("nats://%s", cfg.NatsAddress),
	)

	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)
	defer events.Close()

	router := newRouter()
	log.Println("Started server on port 8080...")
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
    router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)
	return 
}
