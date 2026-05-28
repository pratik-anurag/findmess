package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/findmesh/findmesh/backend/internal/config"
	fmcrypto "github.com/findmesh/findmesh/backend/internal/crypto"
	"github.com/findmesh/findmesh/backend/internal/db"
)

type Service struct {
	store *db.Store
	cfg   config.Config
}

func NewService(store *db.Store, cfg config.Config) *Service {
	return &Service{store: store, cfg: cfg}
}

func (s *Service) StartOTP(phone string) (string, error) {
	phone = fmcrypto.NormalizePhone(phone)
	if phone == "" {
		return "", errors.New("phone is required")
	}
	hash := fmcrypto.HashPhone(phone, s.cfg.PhonePepper)
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	s.store.OTPs[hash] = s.cfg.DevOTP
	return s.cfg.DevOTP, nil
}

func (s *Service) VerifyOTP(phone, otp string) (*db.User, string, error) {
	phone = fmcrypto.NormalizePhone(phone)
	hash := fmcrypto.HashPhone(phone, s.cfg.PhonePepper)
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()

	expected := s.store.OTPs[hash]
	if expected == "" || expected != otp {
		return nil, "", errors.New("invalid otp")
	}
	now := time.Now().UTC()
	userID := s.store.UsersByPhoneHash[hash]
	var user *db.User
	if userID == "" {
		encrypted, err := fmcrypto.EncryptString(phone, s.cfg.StorageKey)
		if err != nil {
			return nil, "", err
		}
		user = &db.User{
			ID:             db.NewID(),
			PhoneHash:      hash,
			PhoneEncrypted: encrypted,
			Status:         "active",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		s.store.Users[user.ID] = user
		s.store.UsersByPhoneHash[hash] = user.ID
	} else {
		user = s.store.Users[userID]
		user.Status = "active"
		user.UpdatedAt = now
		user.DeletedAt = nil
	}
	delete(s.store.OTPs, hash)
	token := fmcrypto.RandomToken(32)
	s.store.Sessions[token] = &db.Session{
		Token:     token,
		UserID:    user.ID,
		Role:      "user",
		CreatedAt: now,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
	}
	return user, token, nil
}

func (s *Service) Logout(token string) {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	delete(s.store.Sessions, token)
}

func (s *Service) DeleteAccount(userID string) error {
	s.store.Mu.Lock()
	defer s.store.Mu.Unlock()
	user := s.store.Users[userID]
	if user == nil {
		return errors.New("user not found")
	}
	now := time.Now().UTC()
	user.Status = "deleted"
	user.DeletedAt = &now
	user.UpdatedAt = now
	for _, tag := range s.store.Tags {
		if tag.OwnerUserID == userID {
			tag.OwnerUserID = ""
			tag.Status = "unpaired"
			tag.TagSecretEncrypted = ""
			tag.UpdatedAt = now
		}
	}
	for token, session := range s.store.Sessions {
		if session.UserID == userID {
			delete(s.store.Sessions, token)
		}
	}
	return nil
}

type Principal struct {
	UserID string
	Role   string
}

type contextKey string

const principalKey contextKey = "principal"

func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(principalKey).(Principal)
	return p, ok
}

func UserID(ctx context.Context) string {
	p, _ := PrincipalFromContext(ctx)
	return p.UserID
}

func BearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return ""
	}
	return strings.TrimSpace(header[len("Bearer "):])
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := BearerToken(r)
		if token == "" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		if token == s.cfg.AdminToken {
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), Principal{UserID: "admin", Role: "admin"})))
			return
		}
		s.store.Mu.RLock()
		session := s.store.Sessions[token]
		s.store.Mu.RUnlock()
		if session == nil || time.Now().UTC().After(session.ExpiresAt) {
			http.Error(w, "invalid bearer token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), Principal{UserID: session.UserID, Role: session.Role})))
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := PrincipalFromContext(r.Context())
		if !ok || p.Role != "admin" {
			http.Error(w, "admin role required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
