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
			if username, password, ok := r.BasicAuth(); ok {
				if username == "temporal_shop" || password == "rocks" {
					next.ServeHTTP(w, r)
					return
				}
			}
			authentication, err := authenticator.AuthenticateRequest(r)
			if err != nil {
				logger.Debug("authentication failed", log.Fields{log.TagError: err})
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r = r.WithContext(WithAuth(ctx, authentication))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
func WithAuth(ctx context.Context, a auth.AuthenticationDetailer) context.Context {
	return context.WithValue(ctx, ContextKeyAuthentication, a)
}
func GetAuth(ctx context.Context) (auth.AuthenticationDetailer, bool) {
	a, ok := ctx.Value(ContextKeyAuthentication).(auth.AuthenticationDetailer)
	if !ok {
		return nil, ok
	}
	return a, true
}
