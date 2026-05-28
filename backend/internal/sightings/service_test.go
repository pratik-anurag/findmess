package sightings

import (
	"crypto/ed25519"
	"testing"
	"time"

	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/protocol"
)

func TestScoring(t *testing.T) {
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	got := Score(ScoreInput{
		SourceType:           protocol.SourceMerchantStand,
		IndependentSightings: 2,
		SeenAt:               now,
		RSSIBucket:           protocol.RSSINear,
		SourceReputationHigh: true,
		MatchingZoneBeacon:   true,
		UserAppCorroboration: true,
	}, now)
	if got != 100 {
		t.Fatalf("score got %d want capped 100", got)
	}
	low := Score(ScoreInput{
		SourceType:               protocol.SourceUserApp,
		SeenAt:                   now.Add(-2 * time.Hour),
		RSSIBucket:               protocol.RSSIFar,
		DuplicateReplaySuspected: true,
	}, now)
	if low != 0 {
		t.Fatalf("low score got %d want 0", low)
	}
}

func TestIngestSignedStandSightingMatchesLostModeAndDeduplicates(t *testing.T) {
	store := db.NewMemoryStore()
	now := time.Date(2026, 5, 28, 10, 0, 0, 0, time.UTC)
	svc := NewService(store)
	svc.Now = func() time.Time { return now }
	pub, priv := fmcrypto.GenerateEd25519KeyPair()
	secret := []byte("0123456789abcdef0123456789abcdef")
	bucket := protocol.BucketTime(now)
	eph := fmcrypto.DeriveEphemeralID(secret, protocol.EpochForTime(bucket))

	store.Mu.Lock()
	user := &db.User{ID: db.NewID(), PhoneHash: "h", Status: "active", CreatedAt: now, UpdatedAt: now}
	merchant := &db.Merchant{ID: db.NewID(), Name: "m", DisplayName: "m", Status: "verified", RecoveryEnabled: true, CreatedAt: now, UpdatedAt: now}
	zone := &db.MerchantZone{ID: db.NewID(), MerchantID: merchant.ID, DisplayArea: "near a participating merchant zone", LocationPrecisionMeters: 500, PublicVisibility: "coarse_only", CreatedAt: now}
	stand := &db.Stand{ID: db.NewID(), MerchantID: merchant.ID, ZoneID: zone.ID, PublicKey: pub, Status: "online", LastHeartbeatAt: &now, CreatedAt: now, UpdatedAt: now}
	tag := &db.Tag{ID: db.NewID(), OwnerUserID: user.ID, Status: "lost", PublicLabel: "bag", TagSecretEncrypted: fmcrypto.Base64Secret(secret), CreatedAt: now, UpdatedAt: now}
	lost := &db.LostModeSession{ID: db.NewID(), TagID: tag.ID, OwnerUserID: user.ID, Status: "active", CreatedAt: now}
	store.Users[user.ID] = user
	store.Merchants[merchant.ID] = merchant
	store.MerchantZones[zone.ID] = zone
	store.Stands[stand.ID] = stand
	store.Tags[tag.ID] = tag
	store.LostModeSessions[lost.ID] = lost
	store.Mu.Unlock()

	payload := protocol.SightingPayload{
		ProtocolVersion: protocol.Version,
		SourceType:      protocol.SourceMerchantStand,
		SourceID:        stand.ID,
		TagEphemeralID:  eph,
		ZoneID:          zone.ID,
		TimeBucket:      bucket,
		RSSIBucket:      protocol.RSSINear,
		Nonce:           "nonce-1",
	}
	payload.Signature = fmcrypto.SignEd25519(ed25519.PrivateKey(priv), []byte(protocol.CanonicalSightingString(payload)))
	sighting, err := svc.Ingest(payload, "")
	if err != nil {
		t.Fatal(err)
	}
	if sighting.ConfidenceScore < 70 {
		t.Fatalf("expected high confidence, got %d", sighting.ConfidenceScore)
	}
	store.Mu.RLock()
	summaryCount := len(store.LastSeenSummaries)
	store.Mu.RUnlock()
	if summaryCount != 1 {
		t.Fatalf("expected one last seen summary, got %d", summaryCount)
	}
	duplicate, err := svc.Ingest(payload, "")
	if err != nil {
		t.Fatal(err)
	}
	if !duplicate.Suspicious {
		t.Fatal("duplicate sighting should be flagged suspicious")
	}
}
