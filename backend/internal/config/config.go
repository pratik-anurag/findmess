package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env              string
	HTTPAddr         string
	APIBaseURL       string
	PhonePepper      string
	StorageKey       string
	DevOTP           string
	AdminToken       string
	RawRetention     time.Duration
	LostRetention    time.Duration
	UseInMemoryEvent bool
}

func Load() Config {
	return Config{
		Env:              getenv("FINDMESH_ENV", "development"),
		HTTPAddr:         getenv("FINDMESH_HTTP_ADDR", ":8080"),
		APIBaseURL:       getenv("FINDMESH_API_BASE_URL", "http://localhost:8080"),
		PhonePepper:      getenv("FINDMESH_PHONE_PEPPER", "dev-phone-pepper-change-me"),
		StorageKey:       getenv("FINDMESH_STORAGE_KEY", "dev-storage-key-change-me-32-bytes"),
		DevOTP:           getenv("FINDMESH_DEV_OTP", "123456"),
		AdminToken:       getenv("FINDMESH_ADMIN_TOKEN", "dev-admin-token"),
		RawRetention:     durationDays("FINDMESH_RAW_RETENTION_DAYS", 30),
		LostRetention:    durationDays("FINDMESH_LOST_RETENTION_DAYS", 30),
		UseInMemoryEvent: getenv("FINDMESH_EVENT_BUS", "memory") == "memory",
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationDays(key string, fallback int) time.Duration {
	raw := getenv(key, strconv.Itoa(fallback))
	days, err := strconv.Atoi(raw)
	if err != nil || days < 1 {
		days = fallback
	}
	return time.Duration(days) * 24 * time.Hour
}
