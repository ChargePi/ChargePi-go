package tls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

func CreateWssClient(CACertificatePath, ClientCertificatePath, ClientKeyPath string) (*ws.Client, error) {
	log.Debugf("Creating a TLS client")

	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.WithError(err).Error("Cannot fetch certificate")
		return nil, err
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(CACertificatePath)
	if err != nil {
		log.WithError(err).Errorf("Error reading certificate")
		return nil, err
	} else if !certPool.AppendCertsFromPEM(caCert) {
		log.Errorf("no ca.cert file found, will use system CA certificates")
		return nil, err
	}

	// Load client certificate
	certificate, err := tls.LoadX509KeyPair(ClientCertificatePath, ClientKeyPath)
	if err != nil {
		log.WithError(err).Errorf("Couldn't load client TLS certificate")
		return nil, err
	}

	// Create client with TLS config
	return ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{certificate},
	}), err
}
