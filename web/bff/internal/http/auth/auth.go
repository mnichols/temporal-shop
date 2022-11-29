package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"net/http"
	"strings"
	"time"
)

const sessionCookieName = "session"
const headerAuthorization = "authorization"

type Claims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}
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
		if authorization := r.Header.Get(headerAuthorization); authorization != "" {
			bearer := strings.Replace(authorization, "Bearer ", "", -1)
			claims := &Claims{}
			_, err := jwt.ParseWithClaims(bearer, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(key), nil
			})
			if err != nil {
				return "", err
			}
			return claims.Email, nil
		}

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
func generateJWT(key, email string) (*jwt.Token, error) {
	id, err := shopping.GenerateShopperHash(key, email)
	if err != nil {
		return nil, err
	}
	claims := Claims{
		jwt.RegisteredClaims{
			ID:        id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 3)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		email,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken, nil
}
func generateJWTString(key, email string) (string, error) {
	tok, err := generateJWT(key, email)
	if err != nil {
		return "", err
	}
	jwtStr, err := tok.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return jwtStr, nil
}
func tokenizeRequest(key, email string, r *http.Request) error {
	str, err := generateJWTString(key, email)
	if err != nil {
		return err
	}
	if r.Header.Get(headerAuthorization) != "" {
		return fmt.Errorf("authorization already exists")
	}
	r.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %v", str))
	return nil
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
