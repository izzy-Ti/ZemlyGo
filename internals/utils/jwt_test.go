package utils

import (
	"testing"
	"time"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "test-secret-key-12345"
	claims := JWTClaims{
		UserID:         42,
		SupabaseUserID: "sb-user-42",
		Email:          "test@example.com",
		Role:           "driver",
	}

	token, exp, err := GenerateJWT(secret, claims, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Fatal("generated token is empty")
	}

	if exp.Before(time.Now()) {
		t.Fatal("expiration time is in the past")
	}

	// Validate valid token
	validatedClaims, err := ValidateJWT(secret, token)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if validatedClaims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", validatedClaims.UserID)
	}
	if validatedClaims.Email != "test@example.com" {
		t.Errorf("expected Email test@example.com, got %s", validatedClaims.Email)
	}
	if validatedClaims.Role != "driver" {
		t.Errorf("expected Role driver, got %s", validatedClaims.Role)
	}

	// Validate with wrong secret
	_, err = ValidateJWT("wrong-secret-key", token)
	if err == nil {
		t.Fatal("expected validation to fail with wrong secret, but succeeded")
	}

	// Validate expired token
	expiredToken, _, err := GenerateJWT(secret, claims, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired JWT: %v", err)
	}
	_, err = ValidateJWT(secret, expiredToken)
	if err == nil {
		t.Fatal("expected expired token to fail validation, but succeeded")
	}
}
