package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/findmesh/findmesh/backend/internal/config"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/server"
)

func TestAuthPairLostAndAbuseFlow(t *testing.T) {
	app := server.NewApp(config.Load(), db.NewMemoryStore(), nil)
	router := app.Router()

	post(t, router, "", "/v1/auth/otp/start", map[string]any{"phone": "+15551234567"}, http.StatusOK)
	verify := post(t, router, "", "/v1/auth/otp/verify", map[string]any{"phone": "+15551234567", "otp": "123456"}, http.StatusOK)
	token := verify["token"].(string)

	pair := post(t, router, token, "/v1/tags/pair/complete", map[string]any{
		"serial":           "FM-TAG-TEST-1",
		"public_label":     "Backpack",
		"firmware_version": "tag-test",
	}, http.StatusOK)
	tag := pair["tag"].(map[string]any)
	tagID := tag["id"].(string)

	lost := post(t, router, token, "/v1/tags/"+tagID+"/lost-mode", map[string]any{
		"safe_message": "If found, contact via FindMesh.",
	}, http.StatusOK)
	if lost["public_lost_token"] == "" {
		t.Fatal("expected public lost token")
	}

	report := post(t, router, token, "/v1/abuse/reports", map[string]any{
		"tag_id":      tagID,
		"category":    "unknown_tracker_alert",
		"description": "Repeated nearby observations.",
	}, http.StatusCreated)
	if report["status"] != "open" {
		t.Fatalf("expected open abuse report, got %#v", report["status"])
	}
}

func post(t *testing.T, h http.Handler, token, path string, body map[string]any, want int) map[string]any {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != want {
		t.Fatalf("%s got status %d want %d body %s", path, rr.Code, want, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	return out
}
