package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/impez/kora/internal/database"
)

type LoginInput struct {
	Username string
	Password string
}

type Service struct {
	DB        *database.Queries
	JWTSecret string
}

func (s *Service) SignToken(username string) (string, error) {
	return signToken(s.JWTSecret, username)
}

func (s *Service) VerifyToken(tokenStr string) (string, error) {
	return verifyToken(s.JWTSecret, tokenStr)
}

func (s *Service) Authenticate(ctx context.Context, input LoginInput) error {
	user, err := s.DB.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}