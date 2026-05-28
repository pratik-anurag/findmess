package db

import (
	"time"
)

func SeedDemo(store *Store) {
	store.Mu.Lock()
	defer store.Mu.Unlock()
	now := time.Now().UTC()
	user := &User{
		ID:             "00000000-0000-4000-8000-000000000101",
		PhoneHash:      "demo-phone-hash",
		PhoneEncrypted: "encrypted-demo-phone",
		Status:         "active",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	merchant := &Merchant{
		ID:              "00000000-0000-4000-8000-000000000201",
		Name:            "Demo Corner Store",
		DisplayName:     "Demo Store",
		Status:          "verified",
		City:            "Bengaluru",
		Category:        "retail",
		RecoveryEnabled: true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	zone := &MerchantZone{
		ID:                      "00000000-0000-4000-8000-000000000202",
		MerchantID:              merchant.ID,
		CoarseGeohash:           "tdr1w",
		DisplayArea:             "near a participating merchant zone in Indiranagar",
		LocationPrecisionMeters: 500,
		PublicVisibility:        "coarse_only",
		CreatedAt:               now,
	}
	tag := &Tag{
		ID:              "00000000-0000-4000-8000-000000000301",
		SerialHash:      "demo-tag-serial-hash",
		OwnerUserID:     user.ID,
		Status:          "lost",
		PublicLabel:     "Keys",
		FirmwareVersion: "tag-dev",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	stand := &Stand{
		ID:              "00000000-0000-4000-8000-000000000401",
		MerchantID:      merchant.ID,
		ZoneID:          zone.ID,
		SerialHash:      "demo-stand-serial-hash",
		Status:          "online",
		FirmwareVersion: "stand-dev",
		PowerSource:     "usb_c",
		WiFiStatus:      "connected",
		LastHeartbeatAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	lost := &LostModeSession{
		ID:              "00000000-0000-4000-8000-000000000501",
		TagID:           tag.ID,
		OwnerUserID:     user.ID,
		Status:          "active",
		SafeMessage:     "If found, contact me via FindMesh.",
		PublicLostToken: "demo-lost-token",
		CreatedAt:       now,
	}
	release := &FirmwareRelease{
		ID:            "00000000-0000-4000-8000-000000000601",
		DeviceType:    "merchant_stand",
		Version:       "0.1.0",
		ManifestURL:   "https://example.invalid/findmesh/merchant-stand/0.1.0.json",
		BinaryURL:     "https://example.invalid/findmesh/merchant-stand/0.1.0.bin",
		Signature:     "dev-signature",
		RolloutStatus: "staged",
		CreatedAt:     now,
	}
	store.Users[user.ID] = user
	store.UsersByPhoneHash[user.PhoneHash] = user.ID
	store.Merchants[merchant.ID] = merchant
	store.MerchantZones[zone.ID] = zone
	store.Tags[tag.ID] = tag
	store.TagsBySerial[tag.SerialHash] = tag.ID
	store.Stands[stand.ID] = stand
	store.StandsBySerial[stand.SerialHash] = stand.ID
	store.LostModeSessions[lost.ID] = lost
	store.FirmwareReleases[release.ID] = release
}
