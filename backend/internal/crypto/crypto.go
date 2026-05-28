package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

const EphemeralIntervalSeconds int64 = 15 * 60

func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}

func HashValue(value, pepper string) string {
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func HashPhone(phone, pepper string) string {
	return HashValue(NormalizePhone(phone), pepper)
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(fmt.Sprintf("random bytes: %v", err))
	}
	return b
}

func RandomToken(n int) string {
	return base64.RawURLEncoding.EncodeToString(RandomBytes(n))
}

func HKDFSHA256(secret, info []byte, length int) []byte {
	salt := make([]byte, sha256.Size)
	mac := hmac.New(sha256.New, salt)
	mac.Write(secret)
	prk := mac.Sum(nil)

	var okm []byte
	var previous []byte
	counter := byte(1)
	for len(okm) < length {
		mac = hmac.New(sha256.New, prk)
		mac.Write(previous)
		mac.Write(info)
		mac.Write([]byte{counter})
		previous = mac.Sum(nil)
		okm = append(okm, previous...)
		counter++
	}
	return okm[:length]
}

func EpochForUnix(unixSeconds int64) int64 {
	return unixSeconds / EphemeralIntervalSeconds
}

func DeriveEphemeralID(tagSecret []byte, epoch int64) string {
	info := make([]byte, len("findmesh-ephemeral-id")+8)
	copy(info, []byte("findmesh-ephemeral-id"))
	binary.BigEndian.PutUint64(info[len("findmesh-ephemeral-id"):], uint64(epoch))
	return hex.EncodeToString(HKDFSHA256(tagSecret, info, 16))
}

func EncryptString(plaintext, keyMaterial string) (string, error) {
	key := sha256.Sum256([]byte(keyMaterial))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := RandomBytes(gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(encoded, keyMaterial string) (string, error) {
	key := sha256.Sum256([]byte(keyMaterial))
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce := raw[:gcm.NonceSize()]
	body := raw[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, body, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func GenerateEd25519KeyPair() (publicKeyBase64 string, privateKey ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("ed25519 keygen: %v", err))
	}
	return base64.RawStdEncoding.EncodeToString(pub), priv
}

func SignEd25519(privateKey ed25519.PrivateKey, payload []byte) string {
	return base64.RawStdEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
}

func VerifyEd25519(publicKeyBase64, signatureBase64 string, payload []byte) bool {
	pub, err := base64.RawStdEncoding.DecodeString(publicKeyBase64)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return false
	}
	sig, err := base64.RawStdEncoding.DecodeString(signatureBase64)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(ed25519.PublicKey(pub), payload, sig)
}

func Base64Secret(secret []byte) string {
	return base64.RawStdEncoding.EncodeToString(secret)
}

func DecodeBase64Secret(encoded string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(encoded)
}
