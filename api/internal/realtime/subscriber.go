package realtime

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/model"
)

// Subscriber wraps a NATS connection for subscribing to events.
type Subscriber struct {
	conn    *nats.Conn
	logger  *slog.Logger
	subject []string
}

// NewSubscriber creates a new subscriber connection.
func NewSubscriber(cfg config.NATSConfig, logger *slog.Logger) (*Subscriber, error) {
	opts := []nats.Option{
		nats.Name(cfg.ClientID + "-subscriber"),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	return &Subscriber{
		conn:   conn,
		logger: logger,
	}, nil
}

// SubscribeUsersCreated subscribes to users.created events.
func (s *Subscriber) SubscribeUsersCreated(handler func(userID, email string)) error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("subscriber not initialized")
	}
	sub, err := s.conn.Subscribe("users.created", func(msg *nats.Msg) {
		var payload struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		}
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			s.logger.Warn("failed to unmarshal users.created", slog.Any("error", err))
			return
		}
		handler(payload.ID, payload.Email)
	})
	if err != nil {
		return err
	}
	s.subject = append(s.subject, sub.Subject)
	return nil
}

// SubscribeChases subscribes to chase lifecycle events.
func (s *Subscriber) SubscribeChases(handler func(event string, chase *model.Chase)) error {
	if s == nil || s.conn == nil {
		return fmt.Errorf("subscriber not initialized")
	}
	sub, err := s.conn.Subscribe("chases.*", func(msg *nats.Msg) {
		var evt ChaseEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			s.logger.Warn("failed to unmarshal chase event", slog.Any("error", err))
			return
		}
		handler(evt.Event, evt.Chase)
	})
	if err != nil {
		return err
	}
	s.subject = append(s.subject, sub.Subject)
	return nil
}

// Close closes the subscriber connection.
func (s *Subscriber) Close() {
	if s == nil || s.conn == nil {
		return
	}
	s.conn.Close()
}
