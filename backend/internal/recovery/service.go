package recovery

import (
	"errors"
	"time"

	"github.com/findmesh/findmesh/backend/internal/db"
)

type Service struct {
	store *db.Store
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(lostModeSessionID, merchantID, zoneID string) (*db.RecoveryRequest, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	session := s.store.LostModeSessions[lostModeSessionID]
	if session == nil || session.Status != "active" {
		return nil, errors.New("active lost mode session not found")
	}
	if merchantID != "" {
		merchant := s.store.Merchants[merchantID]
		if merchant == nil || !merchant.RecoveryEnabled {
			return nil, errors.New("merchant recovery is not enabled")
		}
	}
	now := time.Now().UTC()
	req := &db.RecoveryRequest{
		ID:                db.NewID(),
		LostModeSessionID: lostModeSessionID,
		MerchantID:        merchantID,
		ZoneID:            zoneID,
		Status:            "requested",
		MaskedThreadID:    db.NewID(),
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	s.store.RecoveryRequests[req.ID] = req
	return req, nil
}

func (s *Service) UpdateStatus(id, status string) (*db.RecoveryRequest, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	req := s.store.RecoveryRequests[id]
	if req == nil {
		return nil, errors.New("recovery request not found")
	}
	req.Status = status
	req.UpdatedAt = time.Now().UTC()
	return req, nil
}

func (s *Service) Message(id, actorType, body string) (*db.RecoveryRequest, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	req := s.store.RecoveryRequests[id]
	if req == nil {
		return nil, errors.New("recovery request not found")
	}
	req.Messages = append(req.Messages, db.Message{
		ID:        db.NewID(),
		ActorType: actorType,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	})
	req.UpdatedAt = time.Now().UTC()
	return req, nil
}

func (s *Service) List() []*db.RecoveryRequest {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.RecoveryRequest
	for _, req := range s.store.RecoveryRequests {
		out = append(out, req)
	}
	return out
}
