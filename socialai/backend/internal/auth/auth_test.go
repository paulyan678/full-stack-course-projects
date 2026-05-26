package auth

import (
	"testing"
	"time"
)

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "correct horse battery staple" || !CheckPassword(hash, "correct horse battery staple") {
		t.Fatal("valid password did not verify safely")
	}
	if CheckPassword(hash, "wrong password") {
		t.Fatal("wrong password verified")
	}
}

func TestPasswordPolicy(t *testing.T) {
	if _, err := HashPassword("too-short"); err == nil {
		t.Fatal("short password was accepted")
	}
}

func TestTokenIssueValidateAndExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	manager := NewTokenManager([]byte("01234567890123456789012345678901"), time.Hour)
	manager.now = func() time.Time { return now }
	token, err := manager.Issue("sky")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := manager.Validate(token)
	if err != nil || claims.Subject != "sky" {
		t.Fatalf("validate token: claims=%+v err=%v", claims, err)
	}
	manager.now = func() time.Time { return now.Add(2 * time.Hour) }
	if _, err := manager.Validate(token); err == nil {
		t.Fatal("expired token was accepted")
	}
	if _, err := manager.Validate(token + "x"); err == nil {
		t.Fatal("tampered token was accepted")
	}
}
