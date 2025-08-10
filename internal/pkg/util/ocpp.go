package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/pkg/errors"
)

// CreateConnectionUrl creates a connection url from the provided settings
func CreateConnectionUrl(connectionSettings settings.ConnectionSettings) string {
	var (
		serverUrl = fmt.Sprintf("ws://%s", connectionSettings.ServerUri)
	)

	// Replace insecure Websockets
	if connectionSettings.TLS.IsEnabled {
		serverUrl = strings.Replace(serverUrl, "ws", "wss", 1)
	}

	return serverUrl
}

// CreateClient creates a Websocket client based on the settings.
func CreateClient(logger *zap.Logger, connectionSettings settings.ConnectionSettings, pingInterval *string) (*ws.Client, error) {
	logger.Debug("Creating a websocket client")

	client := ws.NewClient()
	clientConfig := ws.NewClientTimeoutConfig()

	if pingInterval != nil {
		// Set the ping interval
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", *pingInterval))
		if err == nil {
			clientConfig.PingPeriod = duration
		}
	}

	// Check if the TLS is enabled for the client
	if connectionSettings.TLS.IsEnabled {
		logger.Debug("TLS enabled for the websocket client")

		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, errors.Wrap(err, "Cannot fetch certificate pool")
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(connectionSettings.TLS.CACertificatePath)
		if err != nil {
			return nil, errors.Wrap(err, "error reading CA certificate")
		} else if !certPool.AppendCertsFromPEM(caCert) {
			return nil, errors.Wrap(err, "no ca.cert file found, will use system CA certificates")
		}

		// Load client certificate
		certificate, err := tls.LoadX509KeyPair(connectionSettings.TLS.ClientCertificatePath, connectionSettings.TLS.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't load client TLS certificate")
		}

		// Create client with TLS config
		client = ws.NewTLSClient(&tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{certificate},
		})
	}

	// If HTTP basic auth is provided, set it in the Websocket client
	if stringUtils.IsNoneEmpty(connectionSettings.BasicAuthUsername, connectionSettings.BasicAuthPassword) {
		logger.Debug("Basic auth enabled")
		client.SetBasicAuth(connectionSettings.BasicAuthUsername, connectionSettings.BasicAuthPassword)
	}

	client.SetTimeoutConfig(clientConfig)
	return client, nil
}
