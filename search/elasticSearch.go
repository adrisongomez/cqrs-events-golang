package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/adrisongomez/cqrs-events-golang/models"
	elastic "github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	repo := ElasticSearchRepository{
		client: client,
	}

	return &repo, nil
}

func (e *ElasticSearchRepository) Close() {
	//
}

func (r *ElasticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	body, err := json.Marshal(feed)

	if err != nil {
		return err
	}

	_, err = r.client.Index(
		"feeds",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(feed.Id),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)

	return err
}

func (r *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) (results []models.Feed, err error) {
	var buff = bytes.Buffer{}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzzines":         3,
				"cutoff_frequency": 0.001,
			},
		},
	}

	if err = json.NewEncoder(&buff).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("feeds"),
		r.client.Search.WithBody(&buff),
		r.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}


	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		feed := models.Feed{}
		source := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(marshal, &feed); err == nil {
			results = append(results, feed)
		} else {
            log.Println(err)
        }
	}

	return results, nil
}
