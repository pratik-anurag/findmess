package sightings

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/protocol"
)

type Service struct {
	store *db.Store
	Now   func() time.Time
}

func NewService(store *db.Store) *Service {
	return &Service{store: store, Now: func() time.Time { return time.Now().UTC() }}
}

func (s *Service) Ingest(p protocol.SightingPayload, authenticatedUserID string) (*db.Sighting, error) {
	if err := protocol.ValidateSightingPayload(p); err != nil {
		return nil, err
	}
	p.TagEphemeralID = strings.ToLower(p.TagEphemeralID)
	p.TimeBucket = protocol.BucketTime(p.TimeBucket)
	canonical := protocol.CanonicalSightingString(p)
	hash := sha256.Sum256([]byte(canonical))
	rawHash := hex.EncodeToString(hash[:])
	dedupKey := p.SourceType + ":" + p.SourceID + ":" + p.Nonce + ":" + p.TagEphemeralID + ":" + p.TimeBucket.Format(time.RFC3339)

	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()

	if p.SourceType == protocol.SourceUserApp {
		p.SourceID = authenticatedUserID
		if p.SourceID == "" {
			return nil, errors.New("authenticated user required for user app sighting")
		}
		dedupKey = p.SourceType + ":" + p.SourceID + ":" + p.Nonce + ":" + p.TagEphemeralID + ":" + p.TimeBucket.Format(time.RFC3339)
	}

	if existingID := s.store.SightingsByDedup[dedupKey]; existingID != "" {
		existing := s.store.Sightings[existingID]
		existing.Suspicious = true
		return existing, nil
	}

	if p.SourceType == protocol.SourceMerchantStand {
		stand := s.store.Stands[p.SourceID]
		if stand == nil {
			return nil, errors.New("stand source not found")
		}
		if !fmcrypto.VerifyEd25519(stand.PublicKey, p.Signature, []byte(canonical)) {
			return nil, errors.New("invalid stand signature")
		}
		if p.ZoneID == "" {
			p.ZoneID = stand.ZoneID
		}
	}

	independent := s.independentSightingsLocked(p)
	score := Score(ScoreInput{
		SourceType:           p.SourceType,
		IndependentSightings: independent,
		SeenAt:               p.TimeBucket,
		RSSIBucket:           p.RSSIBucket,
		SourceReputationHigh: p.SourceType == protocol.SourceMerchantStand,
		StandHealthy:         s.standHealthyLocked(p.SourceID),
		MatchingZoneBeacon:   p.ZoneEphemeralID != "",
		UserAppCorroboration: independent >= 2 && p.SourceType == protocol.SourceMerchantStand,
	}, s.Now())
	sighting := &db.Sighting{
		ID:              db.NewID(),
		TagEphemeralID:  p.TagEphemeralID,
		SourceType:      p.SourceType,
		SourceID:        p.SourceID,
		ZoneID:          p.ZoneID,
		TimeBucket:      p.TimeBucket,
		RSSIBucket:      p.RSSIBucket,
		ConfidenceScore: score,
		Nonce:           p.Nonce,
		Signature:       p.Signature,
		RawPayloadHash:  rawHash,
		Suspicious:      false,
		CreatedAt:       s.Now(),
	}
	s.store.Sightings[sighting.ID] = sighting
	s.store.SightingsByDedup[dedupKey] = sighting.ID
	s.matchLostTagsLocked(sighting)
	return sighting, nil
}

func (s *Service) Batch(payloads []protocol.SightingPayload, userID string) ([]*db.Sighting, []string) {
	var sightings []*db.Sighting
	var errs []string
	for _, p := range payloads {
		sighting, err := s.Ingest(p, userID)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		sightings = append(sightings, sighting)
	}
	return sightings, errs
}

func (s *Service) independentSightingsLocked(p protocol.SightingPayload) int {
	sources := map[string]bool{}
	for _, existing := range s.store.Sightings {
		if existing.TagEphemeralID == p.TagEphemeralID &&
			existing.ZoneID == p.ZoneID &&
			existing.TimeBucket.Equal(p.TimeBucket) &&
			existing.SourceID != p.SourceID {
			sources[existing.SourceType+":"+existing.SourceID] = true
		}
	}
	if len(sources) > 0 {
		return len(sources) + 1
	}
	return 1
}

func (s *Service) standHealthyLocked(standID string) bool {
	stand := s.store.Stands[standID]
	if stand == nil || stand.LastHeartbeatAt == nil {
		return false
	}
	return s.Now().Sub(*stand.LastHeartbeatAt) < 24*time.Hour && stand.LastError == ""
}

func (s *Service) matchLostTagsLocked(sighting *db.Sighting) {
	for _, tag := range s.store.Tags {
		if tag.Status != "lost" || tag.TagSecretEncrypted == "" {
			continue
		}
		secret, err := fmcrypto.DecodeBase64Secret(tag.TagSecretEncrypted)
		if err != nil {
			continue
		}
		epoch := protocol.EpochForTime(sighting.TimeBucket)
		for _, candidateEpoch := range []int64{epoch - 1, epoch, epoch + 1} {
			if fmcrypto.DeriveEphemeralID(secret, candidateEpoch) == sighting.TagEphemeralID {
				s.updateLastSeenLocked(tag, sighting)
				break
			}
		}
	}
}

func (s *Service) updateLastSeenLocked(tag *db.Tag, sighting *db.Sighting) {
	now := s.Now()
	tag.LastSeenAt = &sighting.TimeBucket
	tag.UpdatedAt = now
	var activeSessionID string
	for _, session := range s.store.LostModeSessions {
		if session.TagID == tag.ID && session.Status == "active" {
			activeSessionID = session.ID
			break
		}
	}
	displayArea := "coarse participating area"
	if zone := s.store.MerchantZones[sighting.ZoneID]; zone != nil {
		displayArea = zone.DisplayArea
	}
	summary := &db.LastSeenSummary{
		ID:                db.NewID(),
		TagID:             tag.ID,
		LostModeSessionID: activeSessionID,
		ZoneID:            sighting.ZoneID,
		DisplayArea:       displayArea,
		ConfidenceLevel:   protocol.ConfidenceLevel(sighting.ConfidenceScore),
		ConfidenceScore:   sighting.ConfidenceScore,
		LastSeenAt:        sighting.TimeBucket,
		UpdatedAt:         now,
	}
	s.store.LastSeenSummaries[summary.ID] = summary
}

func (s *Service) Retain(rawRetention time.Duration, lostRetention time.Duration) int {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	now := s.Now()
	deleted := 0
	for id, sighting := range s.store.Sightings {
		retention := rawRetention
		for _, session := range s.store.LostModeSessions {
			if session.Status == "active" && session.CreatedAt.Before(sighting.CreatedAt.Add(lostRetention)) {
				retention = lostRetention
				break
			}
		}
		if now.Sub(sighting.CreatedAt) > retention {
			delete(s.store.Sightings, id)
			deleted++
		}
	}
	return deleted
}

func (s *Service) ListDebug() []*db.Sighting {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.Sighting
	for _, sighting := range s.store.Sightings {
		out = append(out, sighting)
	}
	return out
}
