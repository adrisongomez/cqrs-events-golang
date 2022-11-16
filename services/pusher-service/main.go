package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adrisongomez/cqrs-events-golang/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	hub := NewHub()

	if err != nil {
		log.Fatalf("%v", err)
	}

	n, err := events.NewNatsEventStore(
		fmt.Sprintf("nats://%s", cfg.NatsAddress),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreateFeed(handleCreateFeed(hub))

	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(n)

	defer events.Close()

	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)

	log.Println("Started server on port :8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleCreateFeed(hub *Hub) func(events.CreatedFeedMessage) {
	return func(m events.CreatedFeedMessage) {
		hub.Broadcast(newCreatedFeedMessage(m.Id, m.Title, m.Description, m.CreatedAt), nil)
	}
}
