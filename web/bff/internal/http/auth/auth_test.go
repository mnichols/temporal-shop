package auth

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"net/http"
	"testing"
)

func TestExtractEmailFromCookie(t *testing.T) {
	url := "/api"
	email := "mike@example.org"
	key := "doit"
	token, err := shopping.GenerateShopperHash(key, email)
	assert.NoError(t, err)

	sessionCookie := &http.Cookie{
		Name:   sessionCookieName,
		Value:  token,
		MaxAge: 300,
	}
	type args struct {
		r      *http.Request
		cookie *http.Cookie
	}
	cookied, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			"valid session cookie exists", args{cookied, sessionCookie}, email, nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			extract := ExtractEmailFromRequest(key)
			if tt.args.cookie != nil {
				tt.args.r.AddCookie(tt.args.cookie)
			}
			actual, err := extract(tt.args.r)
			if err != nil {
				if tt.err == nil {
					assert.NoError(err)
				} else {
					assert.EqualError(err, tt.err.Error())
				}
			} else {
				assert.Equal(tt.want, actual)
			}
		})
	}
}

func TestExtractEmailFromHeader(t *testing.T) {
	url := "/api"
	email := "mike@example.org"
	key := "doit"

	type args struct {
		r      *http.Request
		cookie *http.Cookie
	}
	headered, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	tok, err := generateJWT(key, email)
	assert.NoError(t, err)
	jwtStr, err := tok.SignedString([]byte(key))
	assert.NoError(t, err)
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
			assert := assert.New(t)
			extract := ExtractEmailFromRequest(key)
			if tt.args.cookie != nil {
				tt.args.r.AddCookie(tt.args.cookie)
			}
			actual, err := extract(tt.args.r)
			if err != nil {
				if tt.err == nil {
					assert.NoError(err)
				} else {
					assert.EqualError(err, tt.err.Error())
				}
			} else {
				assert.Equal(tt.want, actual)
			}

		})
	}
}
