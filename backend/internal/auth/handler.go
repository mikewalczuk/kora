package auth

import (
	"context"
	"net/http"

	"github.com/impez/kora/internal/api"
)

type Handler struct {
	Service *Service
}

func (h *Handler) Login(ctx context.Context, req api.LoginRequestObject) (api.LoginResponseObject, error) {
	user, err := h.Service.Authenticate(ctx, LoginInput{
		Username: req.Body.Username,
		Password: req.Body.Password,
	})
	if err != nil {
		return api.Login401Response{}, nil
	}

	token, err := h.Service.SignToken(req.Body.Username, user.ID.String())
	if err != nil {
		return nil, err
	}

	return login200Response{token: token}, nil
}

func (h *Handler) Logout(ctx context.Context, req api.LogoutRequestObject) (api.LogoutResponseObject, error) {
	return logout204Response{}, nil
}

func (h *Handler) GetMe(ctx context.Context, req api.GetMeRequestObject) (api.GetMeResponseObject, error) {
	r, _ := ctx.Value(RequestKey{}).(*http.Request)
	if r == nil {
		return api.GetMe401Response{}, nil
	}

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return api.GetMe401Response{}, nil
	}

	username, _, err := h.Service.VerifyToken(cookie.Value)
	if err != nil {
		return api.GetMe401Response{}, nil
	}

	return api.GetMe200JSONResponse{Username: username}, nil
}

// login200Response sets the session cookie then writes 200.
type login200Response struct {
	token string
}

func (r login200Response) VisitLoginResponse(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    r.token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(tokenTTL.Seconds()),
	})
	w.WriteHeader(http.StatusOK)
	return nil
}

// logout204Response clears the session cookie then writes 204.
type logout204Response struct{}

func (r logout204Response) VisitLogoutResponse(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RequestKey is the context key for storing *http.Request.
type RequestKey struct{}
