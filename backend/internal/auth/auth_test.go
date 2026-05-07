package auth

import (
	"context"
	"testing"
)

func TestAuthenticate_ValidCredentials(t *testing.T) {
	svc := setupService(t)
	seedUser(t, "alice", "secret123")

	_, err := svc.Authenticate(context.Background(), LoginInput{
		Username: "alice",
		Password: "secret123",
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestAuthenticate_WrongPassword(t *testing.T) {
	svc := setupService(t)
	seedUser(t, "bob", "correcthorse")

	_, err := svc.Authenticate(context.Background(), LoginInput{
		Username: "bob",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}
