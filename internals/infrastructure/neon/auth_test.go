package neon

import (
	"testing"
	"time"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

func TestNeonAuthValidation(t *testing.T) {
	secret := "neon-test-jwt-secret-key-123"
	client := NewClient(&configs.Config{
		JWTSecret: secret,
	})

	// 1. Generate dev/test token
	claims := utils.JWTClaims{
		UserID:         10,
		SupabaseUserID: "neon_usr_1001",
		Email:          "rider@example.com",
		Role:           "rider",
	}

	token, _, err := utils.GenerateJWT(secret, claims, 2*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// 2. Validate token with Neon Auth client
	neonClaims, err := client.ValidateToken(token)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if neonClaims.AuthID != "neon_usr_1001" {
		t.Errorf("expected AuthID neon_usr_1001, got %s", neonClaims.AuthID)
	}
	if neonClaims.Email != "rider@example.com" {
		t.Errorf("expected email rider@example.com, got %s", neonClaims.Email)
	}
	if neonClaims.Role != "rider" {
		t.Errorf("expected role rider, got %s", neonClaims.Role)
	}

	// 3. Test invalid token
	_, err = client.ValidateToken("invalid.token.structure")
	if err == nil {
		t.Fatal("expected invalid token to fail, but succeeded")
	}

	// 4. Test empty token
	_, err = client.ValidateToken("")
	if err == nil {
		t.Fatal("expected empty token to fail, but succeeded")
	}
}
