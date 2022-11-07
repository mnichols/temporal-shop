package app

import (
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/pkg/shopping"
	"github.com/temporalio/temporal-shop/web/bff/internal/clients/temporal"
	"github.com/temporalio/temporal-shop/web/bff/internal/instrumentation/log"
	"net/http"
)

type Router interface {
	Get(string, http.HandlerFunc)
}
type Handlers struct {
	temporal *temporal.Clients
}

func NewHandlers(opts ...Option) (*Handlers, error) {
	h := &Handlers{}
	for _, opt := range opts {
		opt(h)
	}
	if h.temporal == nil {
		return nil, fmt.Errorf("temporal client required and missing")
	}

	return h, nil
}

func (h *Handlers) GET(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received ", r.URL.String())
	ctx := r.Context()
	logger := log.GetLogger(r.Context())
	wid := shopping.ExtractShopperEmail(cfg.EncryptionKey)
	h.temporal.Client.DescribeWorkflowExecution(ctx)
	h.inner.ServeHTTP(w, r)

}
