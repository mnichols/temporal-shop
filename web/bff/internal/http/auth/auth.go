package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/services/go/pkg/session"
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
	Validate(ctx context.Context, id string) error
	Start(ctx context.Context, params *orchestrations.StartSessionRequest) error
}
type Authentication struct {
	Email string
}
type Authenticator struct {
	sessionStore   SessionStore
	encryptionKey  string
	associatedData []byte
}

func NewAuthenticator(encryptionKey string, associatedData []byte, sessionStore SessionStore) (*Authenticator, error) {
	if sessionStore == nil {
		return nil, fmt.Errorf("session store is required")
	}
	return &Authenticator{
		sessionStore:   sessionStore,
		encryptionKey:  encryptionKey,
		associatedData: associatedData,
	}, nil
}

func (a *Authenticator) AuthenticateRequest(r *http.Request) (*Authentication, error) {
	email, err := ExtractEmailFromRequest(a.encryptionKey)(r)
	if err != nil && errors.Is(err, http.ErrNoCookie) {
		return nil, AuthenticationFailedError
	}
	// TODO move this id creation to a higher middleware
	id, err := session.NewID([]byte(a.encryptionKey), a.associatedData, &values.SessionID{Email: email})
	if err != nil {
		return nil, err
	}
	err = a.sessionStore.Validate(r.Context(), id.String())
	if err != nil {
		return nil, AuthenticationFailedError
	}
	return &Authentication{Email: email}, nil
}
func (a *Authenticator) StartSession(ctx context.Context, email string) error {
	id, err := session.NewID([]byte(a.encryptionKey), a.associatedData, &values.SessionID{Email: email})
	if err != nil {
		return err
	}
	return a.sessionStore.Start(ctx, &orchestrations.StartSessionRequest{
		Id:    id.String(),
		Email: email,
	})
}
func (a *Authenticator) GenerateToken(ctx context.Context, email string) (string, error) {
	tok, err := a.generateJWT(ctx, &values.SessionID{Email: email})
	jwtStr, err := tok.SignedString([]byte(a.encryptionKey))
	if err != nil {
		return "", err
	}
	return jwtStr, nil
}
func (a *Authenticator) generateJWT(_ context.Context, sess *values.SessionID) (*jwt.Token, error) {
	id, err := session.NewID([]byte(a.encryptionKey), a.associatedData, sess)
	if err != nil {
		return nil, err
	}
	claims := Claims{
		jwt.RegisteredClaims{
			ID:        id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 3)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		sess.Email,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken, nil
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
		return "", fmt.Errorf("cookie disabled")
		//
		//sessionCookie, err := r.Cookie(sessionCookieName)
		//if err != nil {
		//	return "", err
		//}
		//email, err := shopping.ExtractShopperEmail(key, sessionCookie.Value, nil)
		//if err != nil {
		//	return "", err
		//}
		//return email, nil
	}
}

//func tokenizeRequest(key, email string, r *http.Request) error {
//	str, err := generateJWTString(key, email)
//	if err != nil {
//		return err
//	}
//	if r.Header.Get(headerAuthorization) != "" {
//		return fmt.Errorf("authorization already exists")
//	}
//	r.Header.Set(headerAuthorization, fmt.Sprintf("Bearer %v", str))
//	return nil
//}

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
