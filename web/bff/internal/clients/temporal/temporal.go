package temporal

import (
	temporalClient "github.com/temporalio/temporal-shop/services/go/pkg/clients/temporal"
)

type Clients = temporalClient.Clients
type Config = temporalClient.Config

var NewClients = temporalClient.NewClients
