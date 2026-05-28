package audit

import (
	"time"

	"github.com/findmesh/findmesh/backend/internal/db"
)

type Logger struct {
	store *db.Store
}

func NewLogger(store *db.Store) *Logger {
	return &Logger{store: store}
}

func (l *Logger) Record(actorType, actorID, action, targetType, targetID string, metadata map[string]any) {
	l.store.Mu.Lock()
	defer l.store.Mu.Unlock()
	event := &db.AuditEvent{
		ID:         db.NewID(),
		ActorType:  actorType,
		ActorID:    actorID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		CreatedAt:  time.Now().UTC(),
	}
	l.store.AuditEvents[event.ID] = event
}
