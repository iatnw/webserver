package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "secret123" {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("secret123")

	// correct password
	match, err := CheckPasswordHash("secret123", hash)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}
	if !match {
		t.Fatal("expected password to match")
	}

	// wrong password
	match, err = CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected password to not match")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "supersecret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	// valid token
	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("expected userID %v, got %v", userID, gotID)
	}

	// wrong secret
	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatal("expected error with wrong secret")
	}

	// expired token
	expiredToken, _ := MakeJWT(userID, secret, -time.Second)
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
