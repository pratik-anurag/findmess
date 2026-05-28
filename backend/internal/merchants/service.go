package merchants

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

func (s *Service) Create(name, displayName, city, category, displayArea string) (*db.Merchant, *db.MerchantZone, error) {
	if name == "" {
		return nil, nil, errors.New("name is required")
	}
	now := time.Now().UTC()
	if displayName == "" {
		displayName = name
	}
	if displayArea == "" {
		displayArea = "participating merchant zone"
	}
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	merchant := &db.Merchant{
		ID:              db.NewID(),
		Name:            name,
		DisplayName:     displayName,
		Status:          "pending_verification",
		City:            city,
		Category:        category,
		RecoveryEnabled: false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	zone := &db.MerchantZone{
		ID:                      db.NewID(),
		MerchantID:              merchant.ID,
		CoarseGeohash:           "unknown",
		DisplayArea:             displayArea,
		LocationPrecisionMeters: 500,
		PublicVisibility:        "coarse_only",
		CreatedAt:               now,
	}
	s.store.Merchants[merchant.ID] = merchant
	s.store.MerchantZones[zone.ID] = zone
	return merchant, zone, nil
}

func (s *Service) SetRecovery(merchantID string, enabled bool) (*db.Merchant, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	m := s.store.Merchants[merchantID]
	if m == nil {
		return nil, errors.New("merchant not found")
	}
	m.RecoveryEnabled = enabled
	m.UpdatedAt = time.Now().UTC()
	return m, nil
}

func (s *Service) Patch(merchantID, displayName, city, category string) (*db.Merchant, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	m := s.store.Merchants[merchantID]
	if m == nil {
		return nil, errors.New("merchant not found")
	}
	if displayName != "" {
		m.DisplayName = displayName
	}
	if city != "" {
		m.City = city
	}
	if category != "" {
		m.Category = category
	}
	m.UpdatedAt = time.Now().UTC()
	return m, nil
}

func (s *Service) List() []*db.Merchant {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.Merchant
	for _, merchant := range s.store.Merchants {
		out = append(out, merchant)
	}
	return out
}
