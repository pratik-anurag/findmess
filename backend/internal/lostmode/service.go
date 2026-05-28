package lostmode

import (
	"errors"
	"time"

	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
	"github.com/findmesh/findmesh/backend/internal/db"
)

type Service struct {
	store *db.Store
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}

func (s *Service) Open(ownerUserID, tagID, safeMessage string) (*db.LostModeSession, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	now := time.Now().UTC()
	for _, session := range s.store.LostModeSessions {
		if session.TagID == tagID && session.Status == "active" {
			session.SafeMessage = safeMessage
			tag.Status = "lost"
			tag.UpdatedAt = now
			return session, nil
		}
	}
	session := &db.LostModeSession{
		ID:              db.NewID(),
		TagID:           tagID,
		OwnerUserID:     ownerUserID,
		Status:          "active",
		SafeMessage:     safeMessage,
		PublicLostToken: fmcrypto.RandomToken(24),
		CreatedAt:       now,
	}
	s.store.LostModeSessions[session.ID] = session
	tag.Status = "lost"
	tag.UpdatedAt = now
	return session, nil
}

func (s *Service) Resolve(ownerUserID, tagID string) (*db.LostModeSession, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	now := time.Now().UTC()
	var resolved *db.LostModeSession
	for _, session := range s.store.LostModeSessions {
		if session.TagID == tagID && session.Status == "active" {
			session.Status = "resolved"
			session.ResolvedAt = &now
			resolved = session
		}
	}
	tag.Status = "recovered"
	tag.UpdatedAt = now
	if resolved == nil {
		return nil, errors.New("active lost mode not found")
	}
	return resolved, nil
}

func (s *Service) FindByPublicToken(token string) (*db.LostModeSession, *db.Tag, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	for _, session := range s.store.LostModeSessions {
		if session.PublicLostToken == token && session.Status == "active" {
			return session, s.store.Tags[session.TagID], nil
		}
	}
	return nil, nil, errors.New("lost item token not found")
}
