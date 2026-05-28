package stands

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

func (s *Service) ClaimStart(serial, publicKey string) (string, *db.Stand, error) {
	if serial == "" {
		return "", nil, errors.New("serial is required")
	}
	serialHash := fmcrypto.HashValue(serial, "findmesh-stand-serial")
	now := time.Now().UTC()
	token := fmcrypto.RandomToken(24)
	tokenHash := fmcrypto.HashValue(token, "findmesh-stand-claim")
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	standID := s.store.StandsBySerial[serialHash]
	stand := s.store.Stands[standID]
	if stand == nil {
		stand = &db.Stand{
			ID:              db.NewID(),
			SerialHash:      serialHash,
			PublicKey:       publicKey,
			Status:          "unclaimed",
			FirmwareVersion: "stand-dev",
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		s.store.Stands[stand.ID] = stand
		s.store.StandsBySerial[serialHash] = stand.ID
	} else if publicKey != "" {
		stand.PublicKey = publicKey
		stand.UpdatedAt = now
	}
	claim := &db.StandClaimToken{
		ID:        db.NewID(),
		StandID:   stand.ID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(30 * time.Minute),
	}
	s.store.StandClaimTokens[claim.ID] = claim
	return token, stand, nil
}

func (s *Service) ClaimComplete(standID, token, merchantID, zoneID string) (*db.Stand, error) {
	now := time.Now().UTC()
	tokenHash := fmcrypto.HashValue(token, "findmesh-stand-claim")
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	stand := s.store.Stands[standID]
	if stand == nil {
		return nil, errors.New("stand not found")
	}
	if s.store.Merchants[merchantID] == nil {
		return nil, errors.New("merchant not found")
	}
	if s.store.MerchantZones[zoneID] == nil {
		return nil, errors.New("merchant zone not found")
	}
	var claim *db.StandClaimToken
	for _, c := range s.store.StandClaimTokens {
		if c.StandID == standID && c.TokenHash == tokenHash && c.ClaimedAt == nil && now.Before(c.ExpiresAt) {
			claim = c
			break
		}
	}
	if claim == nil {
		return nil, errors.New("invalid or expired claim token")
	}
	claim.ClaimedAt = &now
	stand.MerchantID = merchantID
	stand.ZoneID = zoneID
	stand.Status = "online"
	stand.UpdatedAt = now
	return stand, nil
}

func (s *Service) Heartbeat(standID string, hb db.DeviceHeartbeat) (*db.Stand, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	stand := s.store.Stands[standID]
	if stand == nil {
		return nil, errors.New("stand not found")
	}
	now := time.Now().UTC()
	hb.ID = db.NewID()
	hb.StandID = standID
	hb.CreatedAt = now
	s.store.DeviceHeartbeats[hb.ID] = &hb
	stand.FirmwareVersion = hb.FirmwareVersion
	stand.BatteryLevel = hb.BatteryLevel
	stand.PowerSource = hb.PowerSource
	stand.WiFiStatus = "connected"
	stand.LastHeartbeatAt = &now
	stand.LastError = hb.LastError
	stand.UpdatedAt = now
	return stand, nil
}

func (s *Service) Get(standID string) (*db.Stand, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	stand := s.store.Stands[standID]
	if stand == nil {
		return nil, errors.New("stand not found")
	}
	return stand, nil
}

func (s *Service) List() []*db.Stand {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.Stand
	for _, stand := range s.store.Stands {
		out = append(out, stand)
	}
	return out
}
