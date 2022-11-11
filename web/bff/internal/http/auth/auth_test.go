package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"net/http"
	"testing"
)

func TestExtractEmail(t *testing.T) {
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
			if tt.err != nil && err != nil {
				assert.EqualError(err, tt.err.Error())
			} else if tt.err != nil {
				assert.NoError(err)
			} else if err != nil {
				assert.Errorf(err, "expected %s", tt.err)
			} else {
				assert.Equal(tt.want, actual)
			}
		})
	}
}
