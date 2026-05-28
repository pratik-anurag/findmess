package firmware

import (
	"errors"
	"time"

	"github.com/findmesh/findmesh/backend/internal/db"
)

type Manifest struct {
	DeviceType string `json:"device_type"`
	Version    string `json:"version"`
	BinaryURL  string `json:"binary_url"`
	Signature  string `json:"signature"`
	Required   bool   `json:"required"`
}

type Service struct {
	store *db.Store
}

func NewService(store *db.Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateRelease(deviceType, version, manifestURL, binaryURL, signature, rolloutStatus string) (*db.FirmwareRelease, error) {
	if deviceType == "" || version == "" {
		return nil, errors.New("device_type and version are required")
	}
	if rolloutStatus == "" {
		rolloutStatus = "staged"
	}
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	release := &db.FirmwareRelease{
		ID:            db.NewID(),
		DeviceType:    deviceType,
		Version:       version,
		ManifestURL:   manifestURL,
		BinaryURL:     binaryURL,
		Signature:     signature,
		RolloutStatus: rolloutStatus,
		CreatedAt:     time.Now().UTC(),
	}
	s.store.FirmwareReleases[release.ID] = release
	return release, nil
}

func (s *Service) Manifest(deviceType string) (Manifest, error) {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var latest *db.FirmwareRelease
	for _, release := range s.store.FirmwareReleases {
		if release.DeviceType == deviceType && release.RolloutStatus != "disabled" {
			if latest == nil || release.CreatedAt.After(latest.CreatedAt) {
				latest = release
			}
		}
	}
	if latest == nil {
		return Manifest{
			DeviceType: deviceType,
			Version:    "0.0.0",
			BinaryURL:  "",
			Signature:  "",
			Required:   false,
		}, nil
	}
	return Manifest{
		DeviceType: latest.DeviceType,
		Version:    latest.Version,
		BinaryURL:  latest.BinaryURL,
		Signature:  latest.Signature,
		Required:   latest.RolloutStatus == "required",
	}, nil
}
