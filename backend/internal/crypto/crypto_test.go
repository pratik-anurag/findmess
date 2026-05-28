package crypto

import (
	"crypto/ed25519"
	"testing"
)

func TestEphemeralIDVector(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	got := DeriveEphemeralID(secret, 123456)
	want := "9e40ef0c677adae9870809e1cd952fc2"
	if got != want {
		t.Fatalf("ephemeral id mismatch: got %s want %s", got, want)
	}
}

func TestEd25519SignVerify(t *testing.T) {
	pub, priv := GenerateEd25519KeyPair()
	payload := []byte("findmesh signed sighting")
	sig := SignEd25519(ed25519.PrivateKey(priv), payload)
	if !VerifyEd25519(pub, sig, payload) {
		t.Fatal("expected signature to verify")
	}
	if VerifyEd25519(pub, sig, []byte("tampered")) {
		t.Fatal("tampered payload verified")
	}
}

func TestEncryptedStorageRoundTrip(t *testing.T) {
	encrypted, err := EncryptString("+15551234567", "test-key")
	if err != nil {
		t.Fatal(err)
	}
	if encrypted == "+15551234567" {
		t.Fatal("phone was stored in plaintext")
	}
	plain, err := DecryptString(encrypted, "test-key")
	if err != nil {
		t.Fatal(err)
	}
	if plain != "+15551234567" {
		t.Fatalf("decrypt mismatch: %s", plain)
	}
}
