package abuse

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

func (s *Service) Report(reporterUserID, tagID, standID, merchantID, category, description string) (*db.AbuseReport, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}
	now := time.Now().UTC()
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	report := &db.AbuseReport{
		ID:             db.NewID(),
		ReporterUserID: reporterUserID,
		TagID:          tagID,
		StandID:        standID,
		MerchantID:     merchantID,
		Category:       category,
		Description:    description,
		Status:         "open",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	s.store.AbuseReports[report.ID] = report
	return report, nil
}

func (s *Service) Action(id, action string) (*db.AbuseReport, error) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	report := s.store.AbuseReports[id]
	if report == nil {
		return nil, errors.New("abuse report not found")
	}
	report.Status = action
	report.UpdatedAt = time.Now().UTC()
	if action == "disable_tag" && report.TagID != "" {
		if tag := s.store.Tags[report.TagID]; tag != nil {
			tag.Status = "disabled"
		}
	}
	if action == "disable_stand" && report.StandID != "" {
		if stand := s.store.Stands[report.StandID]; stand != nil {
			stand.Status = "disabled"
		}
	}
	return report, nil
}

func (s *Service) List() []*db.AbuseReport {
	s.store.Mu.RLock()
	defer s.store.Mu.RUnlock()
	var out []*db.AbuseReport
	for _, report := range s.store.AbuseReports {
		out = append(out, report)
	}
	return out
}
