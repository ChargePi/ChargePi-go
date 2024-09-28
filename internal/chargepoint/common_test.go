package chargepoint

import (
	"testing"

	"github.com/ChargePi/ChargePi-go/pkg/tls"
	"github.com/stretchr/testify/assert"
)

func TestCreateConnectionUrl(t *testing.T) {
	tests := []struct {
		name     string
		settings ConnectionSettings
		want     string
		err      bool
	}{
		{
			name: "Valid URL",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "ws://localhost:8887/ocpp/2.0.1",
				TLS:             tls.TLS{},
			},
			want: "ws://localhost:8887/ocpp/2.0.1",
			err:  false,
		},
		{
			name: "Valid URL without scheme",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "localhost:8887/ocpp/2.0.1",
				TLS:             tls.TLS{},
			},
			want: "ws://localhost:8887/ocpp/2.0.1",
			err:  true,
		},
		{
			name: "Valid URL with TLS",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "wss://localhost:8887/ocpp/2.0.1",
				TLS: tls.TLS{
					IsEnabled: true,
					// todo
				},
			},
			want: "wss://localhost:8887/ocpp/2.0.1",
			err:  false,
		},
		{
			name: "Invalid TLS certs",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "ws://localhost:8887/ocpp/2.0.1",
				TLS: tls.TLS{
					IsEnabled: true,
					// todo
				},
			},
			want: "wss://localhost:8887/ocpp/2.0.1",
			err:  false,
		},
		{
			name: "Invalid scheme",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "http://localhost:8887/ocpp/2.0.1",
				TLS:             tls.TLS{},
			},
			want: "",
			err:  true,
		},
		{
			name: "Includes query parameters",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "ws://localhost:8887/ocpp/2.0.1?param=1",
				TLS:             tls.TLS{},
			},
			want: "",
			err:  true,
		},
		{
			name: "Invalid url",
			settings: ConnectionSettings{
				Id:              "ChargePoint",
				ProtocolVersion: "1.6",
				ServerUri:       "localhost:8887/ocpp/2.0.1???",
				TLS:             tls.TLS{},
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := CreateConnectionUrl(tt.settings)
			if tt.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, url)
		})
	}
}

func TestCreateClient(t *testing.T) {
	// Untestable
	t.Skipf("Skipping test for %s", t.Name())
}
