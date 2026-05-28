package sightings

import (
	"time"

	"github.com/findmesh/findmesh/backend/internal/protocol"
)

type ScoreInput struct {
	SourceType               string
	IndependentSightings     int
	SeenAt                   time.Time
	RSSIBucket               string
	SourceReputationHigh     bool
	StandHealthy             bool
	Suspicious               bool
	MatchingZoneBeacon       bool
	UserAppCorroboration     bool
	SourceFlagged            bool
	DuplicateReplaySuspected bool
}

func Score(input ScoreInput, now time.Time) int {
	score := 0
	switch input.SourceType {
	case protocol.SourceMerchantStand:
		score += 50
	case protocol.SourceUserApp:
		score += 30
	}
	if input.IndependentSightings >= 2 {
		score += 20
	}
	if input.RSSIBucket == protocol.RSSINear {
		score += 10
	}
	if input.SourceReputationHigh || input.StandHealthy {
		score += 10
	}
	if input.MatchingZoneBeacon {
		score += 10
	}
	if input.UserAppCorroboration {
		score += 10
	}
	if now.Sub(input.SeenAt) > time.Hour {
		score -= 20
	}
	if input.DuplicateReplaySuspected {
		score -= 30
	}
	if input.Suspicious || input.SourceFlagged {
		score -= 50
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
