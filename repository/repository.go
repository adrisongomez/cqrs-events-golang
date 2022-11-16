package repository

import (
	"context"

	"github.com/adrisongomez/cqrs-events-golang/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
	ListFeed(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(r Repository) {
    repository = r
}

func Close() {
    repository.Close()
}


func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

func ListFeed(ctx context.Context) ([]*models.Feed, error) {
	return repository.ListFeed(ctx)
}
