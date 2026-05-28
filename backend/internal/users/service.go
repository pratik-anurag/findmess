package users

import (
	"time"

	"github.com/findmesh/findmesh/backend/internal/db"
)

type Service struct {
	store *db.Store
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}

func (s *Service) Me(userID string) (*db.User, bool) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	u := s.store.Users[userID]
	return u, u != nil
}

func (s *Service) UpsertDevice(userID, platform, pushToken, appVersion string, finderEnabled bool) *db.UserDevice {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	now := time.Now().UTC()
	for _, d := range s.store.UserDevices {
		if d.UserID == userID && d.Platform == platform {
			d.PushToken = pushToken
			d.AppVersion = appVersion
			d.FinderParticipationEnabled = finderEnabled
			d.LastSeenAt = now
			return d
		}
	}
	d := &db.UserDevice{
		ID:                         db.NewID(),
		UserID:                     userID,
		Platform:                   platform,
		PushToken:                  pushToken,
		AppVersion:                 appVersion,
		FinderParticipationEnabled: finderEnabled,
		LastSeenAt:                 now,
		CreatedAt:                  now,
	}
	s.store.UserDevices[d.ID] = d
	return d
}

func (s *Service) ExportAccount(userID string) map[string]any {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	tags := []db.Tag{}
	for _, tag := range s.store.Tags {
		if tag.OwnerUserID == userID {
			copy := *tag
			copy.TagSecretEncrypted = ""
			tags = append(tags, copy)
		}
	}
	return map[string]any{
		"user": s.store.Users[userID],
		"tags": tags,
	}
}
