package auth

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/impez/kora/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type LoginInput struct {
	Username string
	Password string
}

type Service struct {
	DB        *database.Queries
	JWTSecret string
}

func (s *Service) Authenticate(ctx context.Context, input LoginInput) (*database.User, error) {
	user, err := s.DB.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}

func (s *Service) SignToken(username, userID string) (string, error) {
	return signToken(s.JWTSecret, username, userID)
}

func (s *Service) VerifyToken(tokenStr string) (username, userID string, err error) {
	return verifyToken(s.JWTSecret, tokenStr)
}

// CurrentUserID extracts the authenticated user's UUID from the session cookie in ctx.
func (s *Service) CurrentUserID(ctx context.Context) (pgtype.UUID, error) {
	r, _ := ctx.Value(RequestKey{}).(*http.Request)
	if r == nil {
		return pgtype.UUID{}, errors.New("not authenticated")
	}
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return pgtype.UUID{}, errors.New("not authenticated")
	}
	_, userIDStr, err := s.VerifyToken(cookie.Value)
	if err != nil {
		return pgtype.UUID{}, errors.New("not authenticated")
	}
	var id pgtype.UUID
	if err := id.Scan(userIDStr); err != nil {
		return pgtype.UUID{}, errors.New("invalid user id in token")
	}
	return id, nil
}
