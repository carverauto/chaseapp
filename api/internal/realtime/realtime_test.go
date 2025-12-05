package realtime

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/require"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/model"
)

func TestPublishAndSubscribeChaseEvents(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("skipping NATS integration test: %v", r)
		}
	}()

	opts := &server.Options{
		Host:   "127.0.0.1",
		Port:   -1,
		NoLog:  true,
		NoSigs: true,
	}
	srv := test.RunServer(opts)
	t.Cleanup(func() {
		srv.Shutdown()
	})

	cfg := config.NATSConfig{
		URL:           srv.ClientURL(),
		ClientID:      "chaseapp-test",
		MaxReconnects: 1,
		ReconnectWait: time.Second,
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	publisher, err := NewPublisher(cfg, logger)
	require.NoError(t, err)
	t.Cleanup(publisher.Close)

	subscriber, err := NewSubscriber(cfg, logger)
	require.NoError(t, err)
	t.Cleanup(subscriber.Close)

	received := make(chan *model.Chase, 1)
	require.NoError(t, subscriber.SubscribeChases(func(event string, chase *model.Chase) {
		if event == SubjectChaseCreated {
			received <- chase
		}
	}))

	now := time.Now().UTC()
	chase := &model.Chase{
		ID:        uuid.New(),
		Title:     "Integration",
		ChaseType: model.ChaseTypeChase,
		Live:      true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	require.NoError(t, publisher.PublishChase(SubjectChaseCreated, chase))

	select {
	case got := <-received:
		require.Equal(t, chase.ID, got.ID)
		require.Equal(t, chase.Title, got.Title)
	case <-time.After(3 * time.Second):
		t.Fatalf("did not receive chase event")
	}
}
