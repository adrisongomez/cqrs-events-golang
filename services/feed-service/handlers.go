package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/adrisongomez/cqrs-events-golang/events"
	"github.com/adrisongomez/cqrs-events-golang/models"
	"github.com/adrisongomez/cqrs-events-golang/repository"
	"github.com/segmentio/ksuid"
)

type CreateFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateFeedRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC()

	id, err := ksuid.NewRandom()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		Id:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	err = repository.InsertFeed(r.Context(), &feed)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = events.PublishCreatedFeed(r.Context(), &feed)

	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
