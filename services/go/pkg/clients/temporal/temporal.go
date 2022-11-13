package temporal

import temporalClient "github.com/temporalio/temporal-shop/services/go/internal/clients/temporal"

type Clients = temporalClient.Clients
type Config = temporalClient.Config
type MockTemporalClient = temporalClient.MockTemporalClient

var NewClients = temporalClient.NewClients
var WithConfig = temporalClient.WithConfig
var WithLogger = temporalClient.WithLogger
var WithOptions = temporalClient.WithOptions
