package tags

import (
	"errors"
	"strings"
	"time"

	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
	"github.com/findmesh/findmesh/backend/internal/db"
)

type Service struct {
	store *db.Store
}

type PairStart struct {
	PairingChallenge string `json:"pairing_challenge"`
	ExpiresAt        string `json:"expires_at"`
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}

func (s *Service) StartPair() PairStart {
	return PairStart{
		PairingChallenge: fmcrypto.RandomToken(24),
		ExpiresAt:        time.Now().UTC().Add(10 * time.Minute).Format(time.RFC3339),
	}
}

func (s *Service) CompletePair(ownerUserID, serial, label, firmwareVersion string) (*db.Tag, string, error) {
	serial = strings.TrimSpace(serial)
	if serial == "" {
		return nil, "", errors.New("serial is required")
	}
	serialHash := fmcrypto.HashValue(serial, "findmesh-tag-serial")
	now := time.Now().UTC()
	secret := fmcrypto.RandomBytes(32)
	encodedSecret := fmcrypto.Base64Secret(secret)

	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	if existingID := s.store.TagsBySerial[serialHash]; existingID != "" {
		tag := s.store.Tags[existingID]
		if tag.OwnerUserID != "" && tag.OwnerUserID != ownerUserID {
			return nil, "", errors.New("tag already paired")
		}
		tag.OwnerUserID = ownerUserID
		tag.PublicLabel = label
		tag.TagSecretEncrypted = encodedSecret
		tag.FirmwareVersion = firmwareVersion
		tag.Status = "active"
		tag.UpdatedAt = now
		return tag, encodedSecret, nil
	}
	tag := &db.Tag{
		ID:                 db.NewID(),
		SerialHash:         serialHash,
		OwnerUserID:        ownerUserID,
		Status:             "active",
		PublicLabel:        label,
		TagSecretEncrypted: encodedSecret,
		FirmwareVersion:    firmwareVersion,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	s.store.Tags[tag.ID] = tag
	s.store.TagsBySerial[serialHash] = tag.ID
	return tag, encodedSecret, nil
}

func (s *Service) List(ownerUserID string) []*db.Tag {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.Tag
	for _, tag := range s.store.Tags {
		if tag.OwnerUserID == ownerUserID {
			out = append(out, tag)
		}
	}
	return out
}

func (s *Service) GetOwned(ownerUserID, tagID string) (*db.Tag, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

func (s *Service) Patch(ownerUserID, tagID, label string) (*db.Tag, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	if label != "" {
		tag.PublicLabel = label
	}
	tag.UpdatedAt = time.Now().UTC()
	return tag, nil
}

func (s *Service) LastSeen(ownerUserID, tagID string) (*db.LastSeenSummary, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	var best *db.LastSeenSummary
	for _, summary := range s.store.LastSeenSummaries {
		if summary.TagID == tagID && (best == nil || summary.LastSeenAt.After(best.LastSeenAt)) {
			best = summary
		}
	}
	if best == nil {
		return nil, errors.New("last seen not available")
	}
	copy := *best
	copy.ConfidenceScore = 0
	return &copy, nil
}

func (s *Service) RingIntent(ownerUserID, tagID string) (map[string]any, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return nil, errors.New("tag not found")
	}
	return map[string]any{
		"tag_id":     tagID,
		"intent":     "ring_when_nearby",
		"expires_at": time.Now().UTC().Add(2 * time.Minute).Format(time.RFC3339),
	}, nil
}

func (s *Service) Delete(ownerUserID, tagID string) error {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	tag := s.store.Tags[tagID]
	if tag == nil || tag.OwnerUserID != ownerUserID {
		return errors.New("tag not found")
	}
	tag.OwnerUserID = ""
	tag.TagSecretEncrypted = ""
	tag.Status = "unpaired"
	tag.UpdatedAt = time.Now().UTC()
	return nil
}
