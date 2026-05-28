package notifications

import (
	"context"
	"log/slog"
)

type Message struct {
	UserID string
	Title  string
	Body   string
	Data   map[string]string
}

type Dispatcher interface {
	Send(ctx context.Context, msg Message) error
}

type LogDispatcher struct {
	Logger *slog.Logger
}

func (d LogDispatcher) Send(ctx context.Context, msg Message) error {
	logger := d.Logger
	if logger == nil {
		logger = slog.Default()
	}
	logger.InfoContext(ctx, "notification dispatched", "user_id", msg.UserID, "title", msg.Title)
	return nil
}
