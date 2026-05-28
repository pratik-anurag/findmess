package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Store struct {
	Mu sync.RWMutex

	Users            map[string]*User
	UsersByPhoneHash map[string]string
	UserDevices      map[string]*UserDevice
	Sessions         map[string]*Session
	OTPs             map[string]string

	Tags         map[string]*Tag
	TagsBySerial map[string]string

	Merchants     map[string]*Merchant
	MerchantZones map[string]*MerchantZone

	Stands           map[string]*Stand
	StandsBySerial   map[string]string
	StandClaimTokens map[string]*StandClaimToken

	Sightings        map[string]*Sighting
	SightingsByDedup map[string]string

	LostModeSessions  map[string]*LostModeSession
	LastSeenSummaries map[string]*LastSeenSummary
	RecoveryRequests  map[string]*RecoveryRequest
	AbuseReports      map[string]*AbuseReport
	AuditEvents       map[string]*AuditEvent
	FirmwareReleases  map[string]*FirmwareRelease
	DeviceHeartbeats  map[string]*DeviceHeartbeat
}

func NewMemoryStore() *Store {
	return &Store{
		Users:             map[string]*User{},
		UsersByPhoneHash:  map[string]string{},
		UserDevices:       map[string]*UserDevice{},
		Sessions:          map[string]*Session{},
		OTPs:              map[string]string{},
		Tags:              map[string]*Tag{},
		TagsBySerial:      map[string]string{},
		Merchants:         map[string]*Merchant{},
		MerchantZones:     map[string]*MerchantZone{},
		Stands:            map[string]*Stand{},
		StandsBySerial:    map[string]string{},
		StandClaimTokens:  map[string]*StandClaimToken{},
		Sightings:         map[string]*Sighting{},
		SightingsByDedup:  map[string]string{},
		LostModeSessions:  map[string]*LostModeSession{},
		LastSeenSummaries: map[string]*LastSeenSummary{},
		RecoveryRequests:  map[string]*RecoveryRequest{},
		AbuseReports:      map[string]*AbuseReport{},
		AuditEvents:       map[string]*AuditEvent{},
		FirmwareReleases:  map[string]*FirmwareRelease{},
		DeviceHeartbeats:  map[string]*DeviceHeartbeat{},
	}
}

func NewID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Sprintf("random id: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16]),
	)
}

type Session struct {
	Token     string
	UserID    string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type User struct {
	ID             string     `json:"id"`
	PhoneHash      string     `json:"phone_hash,omitempty"`
	PhoneEncrypted string     `json:"-"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type UserDevice struct {
	ID                         string    `json:"id"`
	UserID                     string    `json:"user_id"`
	Platform                   string    `json:"platform"`
	PushToken                  string    `json:"push_token,omitempty"`
	AppVersion                 string    `json:"app_version"`
	FinderParticipationEnabled bool      `json:"finder_participation_enabled"`
	LastSeenAt                 time.Time `json:"last_seen_at"`
	CreatedAt                  time.Time `json:"created_at"`
}

type Tag struct {
	ID                 string     `json:"id"`
	SerialHash         string     `json:"-"`
	OwnerUserID        string     `json:"owner_user_id,omitempty"`
	Status             string     `json:"status"`
	PublicLabel        string     `json:"public_label,omitempty"`
	TagSecretEncrypted string     `json:"-"`
	BatteryLevel       *int       `json:"battery_level,omitempty"`
	FirmwareVersion    string     `json:"firmware_version,omitempty"`
	LastSeenAt         *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type TagOwnershipEvent struct {
	ID        string    `json:"id"`
	TagID     string    `json:"tag_id"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
}

type Merchant struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	DisplayName     string    `json:"display_name"`
	Status          string    `json:"status"`
	City            string    `json:"city,omitempty"`
	Category        string    `json:"category,omitempty"`
	RecoveryEnabled bool      `json:"recovery_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MerchantZone struct {
	ID                      string    `json:"id"`
	MerchantID              string    `json:"merchant_id"`
	CoarseGeohash           string    `json:"coarse_geohash"`
	DisplayArea             string    `json:"display_area"`
	Latitude                *float64  `json:"latitude,omitempty"`
	Longitude               *float64  `json:"longitude,omitempty"`
	LocationPrecisionMeters int       `json:"location_precision_meters"`
	PublicVisibility        string    `json:"public_visibility"`
	CreatedAt               time.Time `json:"created_at"`
}

type Stand struct {
	ID              string     `json:"id"`
	MerchantID      string     `json:"merchant_id,omitempty"`
	ZoneID          string     `json:"zone_id,omitempty"`
	SerialHash      string     `json:"-"`
	PublicKey       string     `json:"public_key,omitempty"`
	Status          string     `json:"status"`
	FirmwareVersion string     `json:"firmware_version,omitempty"`
	BatteryLevel    *int       `json:"battery_level,omitempty"`
	PowerSource     string     `json:"power_source,omitempty"`
	WiFiStatus      string     `json:"wifi_status,omitempty"`
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at,omitempty"`
	LastError       string     `json:"last_error,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type StandClaimToken struct {
	ID        string     `json:"id"`
	StandID   string     `json:"stand_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	ClaimedAt *time.Time `json:"claimed_at,omitempty"`
}

type Sighting struct {
	ID              string    `json:"id"`
	TagEphemeralID  string    `json:"tag_ephemeral_id"`
	SourceType      string    `json:"source_type"`
	SourceID        string    `json:"source_id,omitempty"`
	ZoneID          string    `json:"zone_id,omitempty"`
	TimeBucket      time.Time `json:"time_bucket"`
	RSSIBucket      string    `json:"rssi_bucket"`
	ConfidenceScore int       `json:"confidence_score"`
	Nonce           string    `json:"nonce"`
	Signature       string    `json:"signature,omitempty"`
	RawPayloadHash  string    `json:"raw_payload_hash"`
	Suspicious      bool      `json:"suspicious"`
	CreatedAt       time.Time `json:"created_at"`
}

type LostModeSession struct {
	ID              string     `json:"id"`
	TagID           string     `json:"tag_id"`
	OwnerUserID     string     `json:"owner_user_id"`
	Status          string     `json:"status"`
	SafeMessage     string     `json:"safe_message,omitempty"`
	PublicLostToken string     `json:"public_lost_token,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
}

type LastSeenSummary struct {
	ID                string    `json:"id"`
	TagID             string    `json:"tag_id"`
	LostModeSessionID string    `json:"lost_mode_session_id,omitempty"`
	ZoneID            string    `json:"zone_id,omitempty"`
	DisplayArea       string    `json:"display_area"`
	ConfidenceLevel   string    `json:"confidence_level"`
	ConfidenceScore   int       `json:"confidence_score,omitempty"`
	LastSeenAt        time.Time `json:"last_seen_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RecoveryRequest struct {
	ID                string    `json:"id"`
	LostModeSessionID string    `json:"lost_mode_session_id"`
	MerchantID        string    `json:"merchant_id,omitempty"`
	ZoneID            string    `json:"zone_id,omitempty"`
	Status            string    `json:"status"`
	MaskedThreadID    string    `json:"masked_thread_id,omitempty"`
	Messages          []Message `json:"messages,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ActorType string    `json:"actor_type"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type AbuseReport struct {
	ID             string    `json:"id"`
	ReporterUserID string    `json:"reporter_user_id,omitempty"`
	TagID          string    `json:"tag_id,omitempty"`
	StandID        string    `json:"stand_id,omitempty"`
	MerchantID     string    `json:"merchant_id,omitempty"`
	Category       string    `json:"category"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AuditEvent struct {
	ID         string         `json:"id"`
	ActorType  string         `json:"actor_type"`
	ActorID    string         `json:"actor_id,omitempty"`
	Action     string         `json:"action"`
	TargetType string         `json:"target_type"`
	TargetID   string         `json:"target_id,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

type FirmwareRelease struct {
	ID            string    `json:"id"`
	DeviceType    string    `json:"device_type"`
	Version       string    `json:"version"`
	ManifestURL   string    `json:"manifest_url"`
	BinaryURL     string    `json:"binary_url"`
	Signature     string    `json:"signature"`
	RolloutStatus string    `json:"rollout_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type DeviceHeartbeat struct {
	ID              string    `json:"id"`
	StandID         string    `json:"stand_id"`
	FirmwareVersion string    `json:"firmware_version"`
	BatteryLevel    *int      `json:"battery_level,omitempty"`
	PowerSource     string    `json:"power_source"`
	WiFiRSSI        int       `json:"wifi_rssi"`
	BufferCount     int       `json:"buffer_count"`
	UptimeSeconds   int64     `json:"uptime_seconds"`
	LastError       string    `json:"last_error,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
