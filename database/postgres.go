package database

import (
	"context"
	"database/sql"
	"github.com/adrisongomez/cqrs-events-golang/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPosgrestRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	postgres := &PostgresRepository{
		db: db,
	}

	return postgres, nil
}

func (p *PostgresRepository) Close() {
	p.db.Close()
}

func (p *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)",
		feed.Id,
		feed.Title,
		feed.Description,
	)
	return err
}

func (p *PostgresRepository) ListFeed(ctx context.Context) ([]*models.Feed, error) {

	rows, err := p.db.QueryContext(
		ctx,
		"SELECT id, title, description, created_at FROM feeds",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	feeds := []*models.Feed{}
	for rows.Next() {
		currentFeed := models.Feed{}

		if err := rows.Scan(&currentFeed.Id,
			&currentFeed.Title,
			&currentFeed.Description,
			&currentFeed.CreatedAt,
		); err == nil {
			feeds = append(feeds, &currentFeed)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
