package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"github.com/temporalio/temporal-shop/web/bff/internal/instrumentation/log"
	"net/http"
)

const sessionCookieName = "session"

type SessionStore interface {
	Validate(ctx context.Context, value string) error
}
type Authentication struct {
	Email string
}
type Authenticator struct {
	sessionStore  SessionStore
	encryptionKey string
}

func NewAuthenticator(encryptionKey string, sessionStore SessionStore) (*Authenticator, error) {
	if sessionStore == nil {
		return nil, fmt.Errorf("session store is required")
	}
	return &Authenticator{
		sessionStore:  sessionStore,
		encryptionKey: encryptionKey,
	}, nil
}

func (a *Authenticator) AuthenticateRequest(r *http.Request) (*Authentication, error) {
	email, err := ExtractEmailFromRequest(a.encryptionKey)(r)
	if err != nil && errors.Is(err, http.ErrNoCookie) {
		return nil, AuthenticationFailedError
	}
	err = a.sessionStore.Validate(r.Context(), email)
	if err != nil {
		return nil, AuthenticationFailedError
	}
	return &Authentication{Email: email}, nil
}
func ExtractEmailFromRequest(key string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		sessionCookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			return "", err
		}
		email, err := shopping.ExtractShopperEmail(key, sessionCookie.Value)
		if err != nil {
			return "", err
		}
		return email, nil
	}
}

func TryExtractEmailFromRequest(key string) func(r *http.Request) string {
	return func(r *http.Request) string {
		result, err := ExtractEmailFromRequest(key)(r)
		if err != nil {
			logger := log.GetLogger(r.Context())
			logger.Error("failed to extract email", log.Fields{"err": err})
			return ""
		}
		return result
	}

}
