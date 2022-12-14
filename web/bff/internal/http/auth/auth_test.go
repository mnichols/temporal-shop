package auth

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	orchestrations2 "github.com/temporalio/temporal-shop/api/temporal_shop/orchestrations/v1"
	"github.com/temporalio/temporal-shop/api/temporal_shop/values/v1"
	"net/http"
	"testing"
)

type ThrowingSessionStore struct{}

func (t *ThrowingSessionStore) Validate(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (t *ThrowingSessionStore) Start(ctx context.Context, params *orchestrations2.StartSessionRequest) error {
	//TODO implement me
	panic("implement me")
}

//func TestExtractEmailFromCookie(t *testing.T) {
//	url := "/api"
//	email := "mike@example.org"
//	key := "doit"
//	token, err := shopping.GenerateShopperHash(key, email, nil)
//	assert.NoError(t, err)
//
//	sessionCookie := &http.Cookie{
//		Name:   sessionCookieName,
//		Value:  token,
//		MaxAge: 300,
//	}
//	type args struct {
//		r      *http.Request
//		cookie *http.Cookie
//	}
//	cookied, err := http.NewRequest("GET", url, nil)
//	assert.NoError(t, err)
//
//	tests := []struct {
//		name string
//		args args
//		want string
//		err  error
//	}{
//		{
//			"valid session cookie exists", args{cookied, sessionCookie}, email, nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			assert := assert.New(t)
//			extract := ExtractEmailFromRequest(key)
//			if tt.args.cookie != nil {
//				tt.args.r.AddCookie(tt.args.cookie)
//			}
//			actual, err := extract(tt.args.r)
//			if err != nil {
//				if tt.err == nil {
//					assert.NoError(err)
//				} else {
//					assert.EqualError(err, tt.err.Error())
//				}
//			} else {
//				assert.Equal(tt.want, actual)
//			}
//		})
//	}
//}

func TestExtractEmailFromHeader(t *testing.T) {
	url := "/api"
	email := "mike@example.org"
	key := "doit"

	type args struct {
		r      *http.Request
		cookie *http.Cookie
	}
	A := assert.New(t)
	headered, err := http.NewRequest("GET", url, nil)
	A.NoError(err)
	sut, err := NewAuthenticator(key, []byte(key), &ThrowingSessionStore{})
	A.NoError(err)

	tok, err := sut.generateJWT(context.Background(), &values.SessionID{Email: email})
	A.NoError(err)
	jwtStr, err := tok.SignedString([]byte(key))
	A.NoError(err)
	headered.Header.Set("authorization", fmt.Sprintf("Bearer %v", jwtStr))

	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "valid session cookie exists", args: args{
				r:      headered,
				cookie: nil,
			}, want: email,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			A := assert.New(t)
			extract := ExtractEmailFromRequest(key)
			if tt.args.cookie != nil {
				tt.args.r.AddCookie(tt.args.cookie)
			}
			actual, err := extract(tt.args.r)
			if err != nil {
				if tt.err == nil {
					A.NoError(err)
				} else {
					A.EqualError(err, tt.err.Error())
				}
			} else {
				A.Equal(tt.want, actual)
			}

		})
	}
}
