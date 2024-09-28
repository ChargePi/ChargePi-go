package chargepoint

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateConnectionUrl creates a connection url from the provided settings
func CreateConnectionUrl(connectionSettings ConnectionSettings) (string, error) {
	// Check that the server URI is valid
	connectionUrl, err := url.ParseRequestURI(connectionSettings.ServerUri)
	if err != nil {
		return "", errors.Wrap(err, "invalid server URI")
	}

	// Allow only ws or wss schemes
	if !strings.HasPrefix(connectionUrl.Scheme, "ws") {
		return "", errors.New("invalid scheme")
	}

	if connectionUrl.RawQuery != "" {
		return "", errors.New("server URI cannot contain query parameters")
	}

	// Replace insecure Websockets
	if connectionSettings.TLS.IsEnabled {
		connectionUrl.Scheme = "wss"
	}

	return connectionUrl.String(), nil
}

// CreateClient creates a Websocket client based on the settings. Automatically adds TLS if enabled.
func CreateClient(connectionSettings ConnectionSettings, pingInterval *string) (*ws.Client, error) {
	log.Debug("Creating a websocket client")

	client := ws.NewClient()
	clientConfig := ws.NewClientTimeoutConfig()

	// Set the ping interval if provided
	if pingInterval != nil {
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", *pingInterval))
		if err == nil {
			clientConfig.PingPeriod = duration
		}
	}

	// Check if the TLS is enabled for the client
	if connectionSettings.TLS.IsEnabled {
		log.Debug("TLS enabled for the websocket client")

		config, err := connectionSettings.TLS.ToTlsConfig()
		if err != nil {
			return nil, err
		}

		// Create client with TLS config
		client = ws.NewTLSClient(config)
	}

	// If HTTP basic auth is provided, set it in the Websocket client
	if connectionSettings.BasicAuth != nil {
		if stringUtils.IsNoneEmpty(connectionSettings.BasicAuth.Username, connectionSettings.BasicAuth.Password) {
			log.Debug("Basic auth enabled")
			client.SetBasicAuth(connectionSettings.BasicAuth.Username, connectionSettings.BasicAuth.Password)
		}
	}

	client.SetTimeoutConfig(clientConfig)
	return client, nil
}
