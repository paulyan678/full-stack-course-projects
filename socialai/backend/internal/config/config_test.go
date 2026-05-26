package config

import "testing"

func TestProductionRejectsPlaceholderJWTSecrets(t *testing.T) {
	for _, secret := range []string{"", "too-short", "replace-with-at-least-32-random-characters", "change-me", "your-secret-here"} {
		t.Run(secret, func(t *testing.T) {
			t.Setenv("APP_ENV", "production")
			t.Setenv("JWT_SECRET", secret)
			if _, err := Load(); err == nil {
				t.Fatalf("production accepted JWT secret %q", secret)
			}
		})
	}
}

func TestProductionAcceptsRandomLengthJWTSecret(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "nN7V7l4c8V2eqGdfYH6CkX3vXhXjvM9zQ5uL1b8WfP0")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.JWTSecret) < 32 {
		t.Fatal("accepted secret was not preserved")
	}
}

func TestRejectsUnknownEnvironment(t *testing.T) {
	t.Setenv("APP_ENV", "prod-ish")
	if _, err := Load(); err == nil {
		t.Fatal("unknown APP_ENV was accepted")
	}
}
