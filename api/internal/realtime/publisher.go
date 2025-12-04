// Package realtime manages NATS publishing for domain events.
package realtime

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/model"
)

const (
	// NATS subjects for chase lifecycle events.
	SubjectChaseCreated = "chases.created"
	SubjectChaseUpdated = "chases.updated"
	SubjectChaseEnded   = "chases.ended"
	SubjectChaseLive    = "chases.live"
	SubjectChaseDeleted = "chases.deleted"
	SubjectAircraftUpdated = "aircraft.updated"
)

// Publisher wraps a NATS connection for publishing events.
type Publisher struct {
	conn   *nats.Conn
	js     *JetStream
	logger *slog.Logger
}

// ChaseEvent is the payload sent to NATS subscribers.
type ChaseEvent struct {
	Event      string       `json:"event"`
	Chase      *model.Chase `json:"chase"`
	OccurredAt time.Time    `json:"occurred_at"`
}

// NewPublisher creates a NATS connection for publishing events.
func NewPublisher(cfg config.NATSConfig, logger *slog.Logger) (*Publisher, error) {
	opts := []nats.Option{
		nats.Name(cfg.ClientID),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	return &Publisher{
		conn:   conn,
		logger: logger,
	}, nil
}

// PublishChase publishes a chase event to the given subject.
func (p *Publisher) PublishChase(subject string, chase *model.Chase) error {
	if p == nil || p.conn == nil {
		return fmt.Errorf("publisher not initialized")
	}

	payload, err := json.Marshal(ChaseEvent{
		Event:      subject,
		Chase:      chase,
		OccurredAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("marshal chase event: %w", err)
	}

	if p.js != nil {
		if err := p.js.Publish(subject, payload); err != nil {
			return fmt.Errorf("publish chase event js: %w", err)
		}
		return nil
	}

	return p.conn.Publish(subject, payload)
}

// IsConnected reports whether the NATS connection is healthy.
func (p *Publisher) IsConnected() bool {
	return p != nil && p.conn != nil && p.conn.Status() == nats.CONNECTED
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p == nil || p.conn == nil {
		return
	}
	p.conn.Close()
}
