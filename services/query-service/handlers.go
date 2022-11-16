package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrisongomez/cqrs-events-golang/events"
	"github.com/adrisongomez/cqrs-events-golang/models"
	"github.com/adrisongomez/cqrs-events-golang/repository"
	"github.com/adrisongomez/cqrs-events-golang/search"
)

func onCreateFeed(msg events.CreatedFeedMessage) {
	feed := models.Feed{
		Id:          msg.Id,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Println(err)
	}
}

func listFeedHanlder(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var err error

	feeds, err := repository.ListFeed(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	/// url/search?q=<query>
	query := r.URL.Query().Get("q")

	if len(query) == 0 {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	results, err := search.SearchFeed(ctx, query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
