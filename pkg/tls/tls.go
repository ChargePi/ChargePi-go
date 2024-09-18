package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// TLS configuration
type TLS struct {
	// Enabled is a flag to enable or disable the TLS
	IsEnabled bool `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`

	// Path to the CA certificate file
	CACertificatePath string `json:"CACertificatePath,omitempty" yaml:"CACertificatePath" mapstructure:"CACertificatePath"`

	// Path to the certificate file
	ClientCertificatePath string `json:"certificatePath,omitempty" yaml:"certificatePath" mapstructure:"certificatePath"`

	// Path to the private key file
	PrivateKeyPath string `json:"keyPath,omitempty" yaml:"keyPath" mapstructure:"keyPath"`
}

func (t *TLS) ToTlsConfig() (*tls.Config, error) {
	if !t.IsEnabled {
		return nil, errors.New("TLS is disabled")
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	// Load CA cert
	caCert, err := os.ReadFile(t.CACertificatePath)
	if err != nil {
		return nil, err
	} else if !certPool.AppendCertsFromPEM(caCert) {
		return nil, err
	}

	// Load client certificate
	certificate, err := tls.LoadX509KeyPair(t.ClientCertificatePath, t.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{certificate},
	}, nil
}
