package server

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/findmesh/findmesh/backend/internal/abuse"
	"github.com/findmesh/findmesh/backend/internal/audit"
	"github.com/findmesh/findmesh/backend/internal/auth"
	"github.com/findmesh/findmesh/backend/internal/config"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/firmware"
	"github.com/findmesh/findmesh/backend/internal/lostmode"
	"github.com/findmesh/findmesh/backend/internal/merchants"
	"github.com/findmesh/findmesh/backend/internal/metrics"
	"github.com/findmesh/findmesh/backend/internal/protocol"
	"github.com/findmesh/findmesh/backend/internal/recovery"
	"github.com/findmesh/findmesh/backend/internal/sightings"
	"github.com/findmesh/findmesh/backend/internal/stands"
	"github.com/findmesh/findmesh/backend/internal/tags"
	"github.com/findmesh/findmesh/backend/internal/users"
)

type App struct {
	Config    config.Config
	Store     *db.Store
	Auth      *auth.Service
	Users     *users.Service
	Tags      *tags.Service
	LostMode  *lostmode.Service
	Merchants *merchants.Service
	Stands    *stands.Service
	Sightings *sightings.Service
	Recovery  *recovery.Service
	Abuse     *abuse.Service
	Firmware  *firmware.Service
	Audit     *audit.Logger
	Metrics   *metrics.Registry
	Logger    *slog.Logger
}

func NewApp(cfg config.Config, store *db.Store, logger *slog.Logger) *App {
	if logger == nil {
		logger = slog.Default()
	}
	return &App{
		Config:    cfg,
		Store:     store,
		Auth:      auth.NewService(store, cfg),
		Users:     users.NewService(store),
		Tags:      tags.NewService(store),
		LostMode:  lostmode.NewService(store),
		Merchants: merchants.NewService(store),
		Stands:    stands.NewService(store),
		Sightings: sightings.NewService(store),
		Recovery:  recovery.NewService(store),
		Abuse:     abuse.NewService(store),
		Firmware:  firmware.NewService(store),
		Audit:     audit.NewLogger(store),
		Metrics:   metrics.NewRegistry(),
		Logger:    logger,
	}
}

func (a *App) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(a.cors())
	r.Use(a.rateLimit())

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Handle("/metrics", a.Metrics.Handler())

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/otp/start", a.otpStart)
		r.Post("/auth/otp/verify", a.otpVerify)
		r.Get("/found/{token}", a.foundItem)
		r.Post("/found/{token}/report", a.foundReport)
		r.Get("/firmware/{device_type}/manifest", a.firmwareManifest)
		r.Post("/sightings", a.ingestSighting)
		r.Post("/sightings/batch", a.ingestSightingBatch)

		r.Group(func(r chi.Router) {
			r.Use(a.Auth.Middleware)
			r.Post("/auth/logout", a.logout)
			r.Get("/me", a.me)
			r.Delete("/me", a.deleteMe)
			r.Post("/devices", a.upsertDevice)

			r.Post("/tags/pair/start", a.tagPairStart)
			r.Post("/tags/pair/complete", a.tagPairComplete)
			r.Get("/tags", a.tagList)
			r.Get("/tags/{tag_id}", a.tagGet)
			r.Patch("/tags/{tag_id}", a.tagPatch)
			r.Post("/tags/{tag_id}/lost-mode", a.lostModeOpen)
			r.Post("/tags/{tag_id}/lost-mode/resolve", a.lostModeResolve)
			r.Get("/tags/{tag_id}/last-seen", a.tagLastSeen)
			r.Post("/tags/{tag_id}/ring-intent", a.tagRingIntent)
			r.Delete("/tags/{tag_id}", a.tagDelete)

			r.Post("/merchants", a.merchantCreate)
			r.Get("/merchants/me", a.merchantMe)
			r.Patch("/merchants/{merchant_id}", a.merchantPatch)
			r.Post("/merchants/{merchant_id}/recovery/enable", a.merchantRecoveryEnable)
			r.Post("/merchants/{merchant_id}/recovery/disable", a.merchantRecoveryDisable)

			r.Post("/stands/claim/start", a.standClaimStart)
			r.Post("/stands/claim/complete", a.standClaimComplete)
			r.Post("/stands/{stand_id}/heartbeat", a.standHeartbeat)
			r.Get("/stands/{stand_id}", a.standGet)
			r.Patch("/stands/{stand_id}", a.standPatch)
			r.Post("/stands/{stand_id}/provisioning-token", a.standProvisioningToken)
			r.Get("/stands/{stand_id}/firmware", a.standFirmware)

			r.Post("/recovery/requests", a.recoveryCreate)
			r.Get("/recovery/requests", a.recoveryList)
			r.Post("/recovery/requests/{id}/accept", a.recoveryAccept)
			r.Post("/recovery/requests/{id}/reject", a.recoveryReject)
			r.Post("/recovery/requests/{id}/message", a.recoveryMessage)
			r.Post("/recovery/requests/{id}/resolve", a.recoveryResolve)

			r.Post("/abuse/reports", a.abuseReport)
		})

		r.Group(func(r chi.Router) {
			r.Use(a.Auth.Middleware)
			r.Use(auth.RequireAdmin)
			r.Get("/sightings/debug", a.adminSightings)
			r.Get("/admin/users", a.adminUsers)
			r.Get("/admin/tags", a.adminTags)
			r.Get("/admin/merchants", a.adminMerchants)
			r.Get("/admin/stands", a.adminStands)
			r.Get("/admin/sightings", a.adminSightings)
			r.Get("/admin/metrics", a.adminMetrics)
			r.Get("/admin/audit-events", a.adminAudit)
			r.Get("/admin/abuse/reports", a.adminAbuseReports)
			r.Post("/admin/abuse/reports/{id}/action", a.adminAbuseAction)
			r.Post("/admin/firmware/releases", a.adminFirmwareRelease)
		})
	})
	return r
}

func (a *App) cors() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *App) rateLimit() func(http.Handler) http.Handler {
	type bucket struct {
		count int
		reset time.Time
	}
	var mu sync.Mutex
	buckets := map[string]*bucket{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = r.RemoteAddr
			}
			now := time.Now().UTC()
			mu.Lock()
			b := buckets[ip]
			if b == nil || now.After(b.reset) {
				b = &bucket{reset: now.Add(time.Minute)}
				buckets[ip] = b
			}
			b.count++
			limited := b.count > 600
			mu.Unlock()
			if limited {
				writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *App) principalFromRequest(r *http.Request) (auth.Principal, bool) {
	if p, ok := auth.PrincipalFromContext(r.Context()); ok {
		return p, true
	}
	token := auth.BearerToken(r)
	if token == "" {
		return auth.Principal{}, false
	}
	if token == a.Config.AdminToken {
		return auth.Principal{UserID: "admin", Role: "admin"}, true
	}
	a.Store.Mu.RLock()
	session := a.Store.Sessions[token]
	a.Store.Mu.RUnlock()
	if session == nil || time.Now().UTC().After(session.ExpiresAt) {
		return auth.Principal{}, false
	}
	return auth.Principal{UserID: session.UserID, Role: session.Role}, true
}

func (a *App) otpStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	otp, err := a.Auth.StartOTP(req.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Metrics.Inc("otp_started")
	writeJSON(w, http.StatusOK, map[string]string{"status": "otp_sent", "dev_otp": otp})
}

func (a *App) otpVerify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	user, token, err := a.Auth.VerifyOTP(req.Phone, req.OTP)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	a.Metrics.Inc("otp_verified")
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	a.Auth.Logout(auth.BearerToken(r))
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

func (a *App) me(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	user, ok := a.Users.Me(p.UserID)
	if !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *App) deleteMe(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	if err := a.Auth.DeleteAccount(p.UserID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "delete_account", "user", p.UserID, nil)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (a *App) upsertDevice(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Platform                   string `json:"platform"`
		PushToken                  string `json:"push_token"`
		AppVersion                 string `json:"app_version"`
		FinderParticipationEnabled bool   `json:"finder_participation_enabled"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	writeJSON(w, http.StatusOK, a.Users.UpsertDevice(p.UserID, req.Platform, req.PushToken, req.AppVersion, req.FinderParticipationEnabled))
}

func (a *App) tagPairStart(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.Tags.StartPair())
}

func (a *App) tagPairComplete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Serial          string `json:"serial"`
		PublicLabel     string `json:"public_label"`
		FirmwareVersion string `json:"firmware_version"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	tag, secret, err := a.Tags.CompletePair(p.UserID, req.Serial, req.PublicLabel, req.FirmwareVersion)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "pair_tag", "tag", tag.ID, nil)
	writeJSON(w, http.StatusOK, map[string]any{"tag": tag, "tag_secret": secret})
}

func (a *App) tagList(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	writeJSON(w, http.StatusOK, a.Tags.List(p.UserID))
}

func (a *App) tagGet(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	tag, err := a.Tags.GetOwned(p.UserID, chi.URLParam(r, "tag_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tag)
}

func (a *App) tagPatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PublicLabel string `json:"public_label"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	tag, err := a.Tags.Patch(p.UserID, chi.URLParam(r, "tag_id"), req.PublicLabel)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tag)
}

func (a *App) lostModeOpen(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SafeMessage string `json:"safe_message"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	session, err := a.LostMode.Open(p.UserID, chi.URLParam(r, "tag_id"), req.SafeMessage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "open_lost_mode", "tag", session.TagID, nil)
	writeJSON(w, http.StatusOK, session)
}

func (a *App) lostModeResolve(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	session, err := a.LostMode.Resolve(p.UserID, chi.URLParam(r, "tag_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "resolve_lost_mode", "tag", session.TagID, nil)
	writeJSON(w, http.StatusOK, session)
}

func (a *App) tagLastSeen(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	summary, err := a.Tags.LastSeen(p.UserID, chi.URLParam(r, "tag_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (a *App) tagRingIntent(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	intent, err := a.Tags.RingIntent(p.UserID, chi.URLParam(r, "tag_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, intent)
}

func (a *App) tagDelete(w http.ResponseWriter, r *http.Request) {
	p, _ := auth.PrincipalFromContext(r.Context())
	if err := a.Tags.Delete(p.UserID, chi.URLParam(r, "tag_id")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "delete_tag", "tag", chi.URLParam(r, "tag_id"), nil)
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (a *App) merchantCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		City        string `json:"city"`
		Category    string `json:"category"`
		DisplayArea string `json:"display_area"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	merchant, zone, err := a.Merchants.Create(req.Name, req.DisplayName, req.City, req.Category, req.DisplayArea)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	a.Audit.Record("user", p.UserID, "create_merchant", "merchant", merchant.ID, nil)
	writeJSON(w, http.StatusCreated, map[string]any{"merchant": merchant, "zone": zone})
}

func (a *App) merchantMe(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.Merchants.List())
}

func (a *App) merchantPatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DisplayName string `json:"display_name"`
		City        string `json:"city"`
		Category    string `json:"category"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	merchant, err := a.Merchants.Patch(chi.URLParam(r, "merchant_id"), req.DisplayName, req.City, req.Category)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, merchant)
}

func (a *App) merchantRecoveryEnable(w http.ResponseWriter, r *http.Request) {
	merchant, err := a.Merchants.SetRecovery(chi.URLParam(r, "merchant_id"), true)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, merchant)
}

func (a *App) merchantRecoveryDisable(w http.ResponseWriter, r *http.Request) {
	merchant, err := a.Merchants.SetRecovery(chi.URLParam(r, "merchant_id"), false)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, merchant)
}

func (a *App) standClaimStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Serial    string `json:"serial"`
		PublicKey string `json:"public_key"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	token, stand, err := a.Stands.ClaimStart(req.Serial, req.PublicKey)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"claim_token": token, "stand": stand})
}

func (a *App) standClaimComplete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StandID    string `json:"stand_id"`
		Token      string `json:"token"`
		MerchantID string `json:"merchant_id"`
		ZoneID     string `json:"zone_id"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	stand, err := a.Stands.ClaimComplete(req.StandID, req.Token, req.MerchantID, req.ZoneID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	a.Audit.Record("user", p.UserID, "claim_stand", "stand", stand.ID, map[string]any{"merchant_id": req.MerchantID})
	writeJSON(w, http.StatusOK, stand)
}

func (a *App) standHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req db.DeviceHeartbeat
	if !decodeJSON(w, r, &req) {
		return
	}
	stand, err := a.Stands.Heartbeat(chi.URLParam(r, "stand_id"), req)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stand)
}

func (a *App) standGet(w http.ResponseWriter, r *http.Request) {
	stand, err := a.Stands.Get(chi.URLParam(r, "stand_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stand)
}

func (a *App) standPatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Status string `json:"status"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	a.Store.Mu.Lock()
	stand := a.Store.Stands[chi.URLParam(r, "stand_id")]
	if stand != nil && req.Status != "" {
		stand.Status = req.Status
		stand.UpdatedAt = time.Now().UTC()
	}
	a.Store.Mu.Unlock()
	if stand == nil {
		writeError(w, http.StatusNotFound, "stand not found")
		return
	}
	writeJSON(w, http.StatusOK, stand)
}

func (a *App) standProvisioningToken(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"stand_id": chi.URLParam(r, "stand_id"),
		"token":    "prov_" + db.NewID(),
	})
}

func (a *App) standFirmware(w http.ResponseWriter, r *http.Request) {
	stand, err := a.Stands.Get(chi.URLParam(r, "stand_id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	deviceType := "merchant_stand"
	if strings.Contains(stand.FirmwareVersion, "tag") {
		deviceType = "lost_tag"
	}
	manifest, _ := a.Firmware.Manifest(deviceType)
	writeJSON(w, http.StatusOK, manifest)
}

func (a *App) ingestSighting(w http.ResponseWriter, r *http.Request) {
	var p protocol.SightingPayload
	if !decodeJSON(w, r, &p) {
		return
	}
	principal, ok := a.principalFromRequest(r)
	if p.SourceType == protocol.SourceUserApp && !ok {
		writeError(w, http.StatusUnauthorized, "user app sightings require authentication")
		return
	}
	sighting, err := a.Sightings.Ingest(p, principal.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Metrics.Inc("sightings_ingested")
	writeJSON(w, http.StatusAccepted, sighting)
}

func (a *App) ingestSightingBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Sightings []protocol.SightingPayload `json:"sightings"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	principal, ok := a.principalFromRequest(r)
	if !ok {
		for _, p := range req.Sightings {
			if p.SourceType == protocol.SourceUserApp {
				writeError(w, http.StatusUnauthorized, "user app sightings require authentication")
				return
			}
		}
	}
	accepted, errs := a.Sightings.Batch(req.Sightings, principal.UserID)
	writeJSON(w, http.StatusAccepted, map[string]any{"accepted": accepted, "errors": errs})
}

func (a *App) recoveryCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LostModeSessionID string `json:"lost_mode_session_id"`
		MerchantID        string `json:"merchant_id"`
		ZoneID            string `json:"zone_id"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	rr, err := a.Recovery.Create(req.LostModeSessionID, req.MerchantID, req.ZoneID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rr)
}

func (a *App) recoveryList(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, a.Recovery.List())
}

func (a *App) recoveryAccept(w http.ResponseWriter, r *http.Request) {
	a.recoveryStatus(w, r, "accepted")
}

func (a *App) recoveryReject(w http.ResponseWriter, r *http.Request) {
	a.recoveryStatus(w, r, "rejected")
}

func (a *App) recoveryResolve(w http.ResponseWriter, r *http.Request) {
	a.recoveryStatus(w, r, "resolved")
}

func (a *App) recoveryStatus(w http.ResponseWriter, r *http.Request, status string) {
	rr, err := a.Recovery.UpdateStatus(chi.URLParam(r, "id"), status)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rr)
}

func (a *App) recoveryMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Body string `json:"body"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	rr, err := a.Recovery.Message(chi.URLParam(r, "id"), "masked_participant", req.Body)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rr)
}

func (a *App) foundItem(w http.ResponseWriter, r *http.Request) {
	session, tag, err := a.LostMode.FindByPublicToken(chi.URLParam(r, "token"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"safe_message": session.SafeMessage,
		"tag_label":    tag.PublicLabel,
		"contact":      "Use FindMesh to send an anonymous found-item report.",
	})
}

func (a *App) foundReport(w http.ResponseWriter, r *http.Request) {
	session, _, err := a.LostMode.FindByPublicToken(chi.URLParam(r, "token"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	rr, err := a.Recovery.Create(session.ID, "", "")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rr)
}

func (a *App) abuseReport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TagID       string `json:"tag_id"`
		StandID     string `json:"stand_id"`
		MerchantID  string `json:"merchant_id"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	p, _ := auth.PrincipalFromContext(r.Context())
	report, err := a.Abuse.Report(p.UserID, req.TagID, req.StandID, req.MerchantID, req.Category, req.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.Audit.Record("user", p.UserID, "create_abuse_report", "abuse_report", report.ID, map[string]any{"category": req.Category})
	writeJSON(w, http.StatusCreated, report)
}

func (a *App) firmwareManifest(w http.ResponseWriter, r *http.Request) {
	manifest, err := a.Firmware.Manifest(chi.URLParam(r, "device_type"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, manifest)
}

func (a *App) adminUsers(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_users", "user", "")
	a.Store.Mu.RLock()
	defer a.Store.Mu.RUnlock()
	var out []*db.User
	for _, user := range a.Store.Users {
		out = append(out, user)
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *App) adminTags(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_tags", "tag", "")
	a.Store.Mu.RLock()
	defer a.Store.Mu.RUnlock()
	var out []*db.Tag
	for _, tag := range a.Store.Tags {
		out = append(out, tag)
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *App) adminMerchants(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_merchants", "merchant", "")
	writeJSON(w, http.StatusOK, a.Merchants.List())
}

func (a *App) adminStands(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_stands", "stand", "")
	writeJSON(w, http.StatusOK, a.Stands.List())
}

func (a *App) adminSightings(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_sightings", "sighting", "")
	writeJSON(w, http.StatusOK, a.Sightings.ListDebug())
}

func (a *App) adminMetrics(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "read_metrics", "metrics", "")
	writeJSON(w, http.StatusOK, map[string]string{"prometheus": "/metrics"})
}

func (a *App) adminAudit(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_audit_events", "audit_event", "")
	a.Store.Mu.RLock()
	defer a.Store.Mu.RUnlock()
	var out []*db.AuditEvent
	for _, event := range a.Store.AuditEvents {
		out = append(out, event)
	}
	writeJSON(w, http.StatusOK, out)
}

func (a *App) adminAbuseReports(w http.ResponseWriter, r *http.Request) {
	a.recordAdminAudit(r, "list_abuse_reports", "abuse_report", "")
	writeJSON(w, http.StatusOK, a.Abuse.List())
}

func (a *App) adminAbuseAction(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	report, err := a.Abuse.Action(chi.URLParam(r, "id"), req.Action)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	a.recordAdminAudit(r, "abuse_action_"+req.Action, "abuse_report", report.ID)
	writeJSON(w, http.StatusOK, report)
}

func (a *App) adminFirmwareRelease(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceType    string `json:"device_type"`
		Version       string `json:"version"`
		ManifestURL   string `json:"manifest_url"`
		BinaryURL     string `json:"binary_url"`
		Signature     string `json:"signature"`
		RolloutStatus string `json:"rollout_status"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	release, err := a.Firmware.CreateRelease(req.DeviceType, req.Version, req.ManifestURL, req.BinaryURL, req.Signature, req.RolloutStatus)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a.recordAdminAudit(r, "create_firmware_release", "firmware_release", release.ID)
	writeJSON(w, http.StatusCreated, release)
}

func (a *App) recordAdminAudit(r *http.Request, action, targetType, targetID string) {
	p, _ := auth.PrincipalFromContext(r.Context())
	a.Audit.Record("admin", p.UserID, action, targetType, targetID, map[string]any{"path": r.URL.Path})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	if msg == "" {
		msg = http.StatusText(status)
	}
	writeJSON(w, status, map[string]string{"error": msg})
}
