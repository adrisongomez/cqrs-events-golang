package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/adrisongomez/cqrs-events-golang/models"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreatedFeedMessage
}

func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)

	if err != nil {
		return nil, err
	}

	store := &NatsEventStore{
		conn: conn,
	}

	return store, nil
}

func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}

	if n.feedCreatedChan != nil {
		close(n.feedCreatedChan)
	}
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}

	err := gob.NewEncoder(&b).Encode(m)

	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(msg)

	if err != err {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

func (n *NatsEventStore) OnCreatedFeed(f func(CreatedFeedMessage)) (err error) {
	msg := CreatedFeedMessage{}

	n.feedCreatedSub, err = n.conn.Subscribe(
		msg.Type(),
		func(m *nats.Msg) {
			n.decodeMessage(m.Data, msg)
			f(msg)
		},
	)

	return

}

func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	m := CreatedFeedMessage{}
	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)

	var err error

	n.feedCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				err := n.decodeMessage(msg.Data, m)
				if err != nil {
					continue
				}
				n.feedCreatedChan <- m
			}
		}
	}()

	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}