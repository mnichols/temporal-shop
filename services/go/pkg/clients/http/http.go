package http

import "github.com/temporalio/temporal-shop/services/go/internal/clients/http"

type Doer = http.Doer
type Config = http.Config
type RequestParams = http.RequestParams

var NewClient = http.NewClient
var MustNewClient = http.MustNewClient
var NewRequest = http.NewRequest
