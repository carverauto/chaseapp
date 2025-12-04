package realtime

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"chaseapp.tv/api/internal/config"
)

// JetStream wraps a JetStream context for publishing and subscribing with durability.
type JetStream struct {
	conn *nats.Conn
	js   nats.JetStreamContext
	log  *slog.Logger
}

// NewJetStream creates a JetStream client.
func NewJetStream(cfg config.NATSConfig, logger *slog.Logger) (*JetStream, error) {
	opts := []nats.Option{
		nats.Name(cfg.ClientID + "-jetstream"),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("init jetstream: %w", err)
	}

	return &JetStream{
		conn: conn,
		js:   js,
		log:  logger,
	}, nil
}

// EnsureStreams creates needed streams if missing.
func (j *JetStream) EnsureStreams(ctx context.Context) error {
	streams := []struct {
		Name     string
		Subjects []string
	}{
		{Name: "chases", Subjects: []string{"chases.*"}},
		{Name: "users", Subjects: []string{"users.*"}},
		{Name: "aircraft", Subjects: []string{"aircraft.*"}},
	}

	for _, s := range streams {
		_, err := j.js.StreamInfo(s.Name, nats.Context(ctx))
		if err == nats.ErrStreamNotFound {
			_, err = j.js.AddStream(&nats.StreamConfig{
				Name:     s.Name,
				Subjects: s.Subjects,
				Storage:  nats.FileStorage,
				MaxAge:   24 * time.Hour,
			}, nats.Context(ctx))
		}
		if err != nil {
			return fmt.Errorf("ensure stream %s: %w", s.Name, err)
		}
	}
	return nil
}

// Subscribe durable consumer with manual ack.
func (j *JetStream) Subscribe(subject, durable string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return j.js.Subscribe(subject, handler, nats.Durable(durable), nats.ManualAck())
}

// Publish publishes to JetStream.
func (j *JetStream) Publish(subject string, data []byte) error {
	_, err := j.js.Publish(subject, data)
	return err
}

// Close closes connections.
func (j *JetStream) Close() {
	if j.conn != nil {
		j.conn.Close()
	}
}
