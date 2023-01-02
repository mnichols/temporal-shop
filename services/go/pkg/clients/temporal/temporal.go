package temporal

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	sdkclient "go.temporal.io/sdk/client"
	"logur.dev/logur"
)

type Clients struct {
	Client          sdkclient.Client
	NamespaceClient sdkclient.NamespaceClient
	Config          *Config
	ClientOptions   sdkclient.Options
	logger          logur.Logger
}

func (c *Clients) Close() error {
	if c.Client != nil {
		c.Client.Close()
	}
	if c.NamespaceClient != nil {
		c.NamespaceClient.Close()
	}
	return nil
}

// NewClients creates the temporal and temporalNamespace clients
func NewClients(ctx context.Context, opts ...Option) (*Clients, error) {
	result := &Clients{
		ClientOptions: sdkclient.Options{},
	}
	opts = append([]Option{
		WithLogger(log.GetLogger(ctx)),
	}, opts...)
	for _, o := range opts {
		o(result)
	}
	logger := log.GetLogger(ctx)
	logger.Info("TEMPORAL CONFIG", map[string]interface{}{
		"cfg": result.Config,
	})

	var cert tls.Certificate
	if result.ClientOptions.Logger == nil && result.logger != nil {
		result.ClientOptions.Logger = logur.LoggerToKV(result.logger)
	}
	// map
	if result.ClientOptions.HostPort == "" {
		result.ClientOptions.HostPort = result.Config.HostPort
	}
	if result.ClientOptions.Namespace == "" {
		result.ClientOptions.Namespace = result.Config.Namespace
	}
	if result.Config.CertFilePath != "" && result.Config.KeyFilePath != "" {
		var err error
		logger.Info("configuring tls", map[string]interface{}{
			"certpath": result.Config.CertFilePath,
			"keypath":  result.Config.KeyFilePath,
		})
		cert, err = tls.LoadX509KeyPair(result.Config.CertFilePath, result.Config.KeyFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS from files: %w", err)
		}
	} else if result.Config.CloudCertPem != "" && result.Config.CloudCertKey != "" {
		var err error
		logger.Info("configuring tls using pem data")
		cert, err = tls.X509KeyPair([]byte(result.Config.CloudCertPem), []byte(result.Config.CloudCertKey))
		if err != nil {
			return nil, fmt.Errorf("failed to load Cloud TLS from data: %w", err)
		}
	}

	if len(cert.Certificate) > 0 {
		result.ClientOptions.ConnectionOptions.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	}
	// validate
	if result.Config.HostPort == "" {
		return nil, fmt.Errorf("temporal HostPort must be defined")
	}
	if result.Config.Namespace == "" {
		return nil, fmt.Errorf("temporal Namespace must be defined")
	}

	logger.Info("creating temporal clients", log.Fields{"client_options_namespace": result.ClientOptions.Namespace})
	var err error
	if result.Client, err = sdkclient.Dial(result.ClientOptions); err != nil {
		return nil, fmt.Errorf("failed to new temporal client %w", err)
	}
	//if result.NamespaceClient, err = sdkclient.NewNamespaceClient(result.ClientOptions); err != nil {
	//	return nil, fmt.Errorf("failed to new namespace client %w", err)
	//}

	return result, nil
}
