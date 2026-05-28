package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
)

const Version = 1

const (
	SourceMerchantStand = "merchant_stand"
	SourceUserApp       = "user_app"
)

const (
	RSSINear   = "near"
	RSSIMedium = "medium"
	RSSIFar    = "far"
)

type LostTagAdvertisement struct {
	AdvType         string            `json:"adv_type"`
	ProtocolVersion int               `json:"protocol_version"`
	EphemeralID     string            `json:"ephemeral_id"`
	Flags           map[string]bool   `json:"flags"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type MerchantZoneAdvertisement struct {
	AdvType           string   `json:"adv_type"`
	ProtocolVersion   int      `json:"protocol_version"`
	ZoneEphemeralID   string   `json:"zone_ephemeral_id"`
	StandCapabilities []string `json:"stand_capabilities"`
}

type SightingPayload struct {
	ProtocolVersion int       `json:"protocol_version"`
	SourceType      string    `json:"source_type"`
	SourceID        string    `json:"source_id,omitempty"`
	TagEphemeralID  string    `json:"tag_ephemeral_id"`
	ZoneEphemeralID string    `json:"zone_ephemeral_id,omitempty"`
	ZoneID          string    `json:"zone_id,omitempty"`
	TimeBucket      time.Time `json:"time_bucket"`
	RSSIBucket      string    `json:"rssi_bucket"`
	Nonce           string    `json:"nonce"`
	Signature       string    `json:"signature,omitempty"`
}

func BucketTime(t time.Time) time.Time {
	return t.UTC().Truncate(time.Duration(fmcrypto.EphemeralIntervalSeconds) * time.Second)
}

func EpochForTime(t time.Time) int64 {
	return fmcrypto.EpochForUnix(BucketTime(t).Unix())
}

func BucketRSSI(rssi int) string {
	switch {
	case rssi >= -60:
		return RSSINear
	case rssi >= -78:
		return RSSIMedium
	default:
		return RSSIFar
	}
}

var ephRe = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

func ValidateSightingPayload(p SightingPayload) error {
	if p.ProtocolVersion != Version {
		return fmt.Errorf("unsupported protocol version %d", p.ProtocolVersion)
	}
	if p.SourceType != SourceMerchantStand && p.SourceType != SourceUserApp {
		return errors.New("invalid source_type")
	}
	if !ephRe.MatchString(p.TagEphemeralID) {
		return errors.New("tag_ephemeral_id must be 16 bytes hex")
	}
	if p.RSSIBucket != RSSINear && p.RSSIBucket != RSSIMedium && p.RSSIBucket != RSSIFar {
		return errors.New("invalid rssi_bucket")
	}
	if strings.TrimSpace(p.Nonce) == "" {
		return errors.New("nonce is required")
	}
	if p.TimeBucket.IsZero() {
		return errors.New("time_bucket is required")
	}
	if p.SourceType == SourceMerchantStand && strings.TrimSpace(p.Signature) == "" {
		return errors.New("merchant stand sightings require signature")
	}
	return nil
}

func CanonicalSightingString(p SightingPayload) string {
	return fmt.Sprintf("v=%d|source_type=%s|source_id=%s|tag=%s|zone_eph=%s|zone_id=%s|time=%s|rssi=%s|nonce=%s",
		p.ProtocolVersion,
		p.SourceType,
		p.SourceID,
		strings.ToLower(p.TagEphemeralID),
		p.ZoneEphemeralID,
		p.ZoneID,
		BucketTime(p.TimeBucket).Format(time.RFC3339),
		p.RSSIBucket,
		p.Nonce,
	)
}

func CanonicalJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func ConfidenceLevel(score int) string {
	switch {
	case score >= 85:
		return "very_high"
	case score >= 60:
		return "high"
	case score >= 30:
		return "medium"
	default:
		return "low"
	}
}
