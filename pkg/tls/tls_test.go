package tls

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/madflojo/testcerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
		{
			name: "TLS is enabled",
			tls: &TLS{
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
		{
			name: "TLS is enabled with invalid CA certificate",
			tls: &TLS{
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
		{
			name: "TLS is enabled with invalid client certificate",
			tls: &TLS{
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
		{
			name: "TLS is enabled with invalid private key",
			tls: &TLS{
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
		{
			name: "TLS is enabled with valid certificates",
			tls: &TLS{
				IsEnabled:             true,
				CACertificatePath:     "",
				ClientCertificatePath: "",
				PrivateKeyPath:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and write self-signed Certificate and Key to temporary files
			cert, key, err := testcerts.GenerateCertsToTempFile("")
			require.NoError(t, err)

			defer os.Remove(key)
			defer os.Remove(cert)

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
