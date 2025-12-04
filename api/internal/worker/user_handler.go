package worker

import (
	"context"
	"log/slog"

	"chaseapp.tv/api/internal/realtime"
	"chaseapp.tv/api/internal/webhook"
)

// UserEventWorker handles user-created events.
type UserEventWorker struct {
	subscriber *realtime.Subscriber
	discord    *webhook.Client
	logger     *slog.Logger
}

// NewUserEventWorker creates a new user event worker.
func NewUserEventWorker(sub *realtime.Subscriber, discord *webhook.Client, logger *slog.Logger) *UserEventWorker {
	return &UserEventWorker{
		subscriber: sub,
		discord:    discord,
		logger:     logger,
	}
}

// Start subscribes to users.created and posts to Discord if configured.
func (w *UserEventWorker) Start(ctx context.Context) error {
	if w.subscriber == nil {
		return nil
	}
	err := w.subscriber.SubscribeUsersCreated(func(userID, email string) {
		w.logger.Info("users.created event", slog.String("user_id", userID), slog.String("email", email))
		if w.discord != nil {
			_ = w.discord.Send(ctx, webhook.Message{
				Content: "New user registered: " + email,
			})
		}
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
