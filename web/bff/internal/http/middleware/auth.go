package middleware

import (
	"context"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/http/auth"
	"net/http"
)

type ContextKey int

const (
	ContextKeyAuthentication = iota
)

func Authenticate(authenticator *auth.Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			logger := log.GetLogger(ctx)
			authentication, err := authenticator.AuthenticateRequest(r)
			if err != nil {
				logger.Debug("authentication failed", log.Fields{log.TagError: err})
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyAuthentication, authentication)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
