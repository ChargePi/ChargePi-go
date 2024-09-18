package tls

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLS_ToTlsConfig(t *testing.T) {
	tests := []struct {
		name    string
		tls     *TLS
		want    *tls.Config
		wantErr bool
	}{
		{
			name: "TLS is disabled",
			tls: &TLS{
				IsEnabled: true,
			},
		},
		{
			name: "TLS is enabled",
		},
		{
			name: "TLS is enabled with invalid CA certificate",
		},
		{
			name: "TLS is enabled with invalid client certificate",
		},
		{
			name: "TLS is enabled with invalid private key",
		},
		{
			name: "TLS is enabled with valid certificates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tls.ToTlsConfig()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
