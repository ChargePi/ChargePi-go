package tls

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func GetTLSClient(CACertificatePath, ClientCertificatePath, ClientKeyPath string) *ws.Client {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.WithError(err).Fatal("Cannot fetch certificate")
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(CACertificatePath)
	if err != nil {
		log.WithError(err).Errorf("Error reading certificate")
		return nil
	} else if !certPool.AppendCertsFromPEM(caCert) {
		log.Errorf("no ca.cert file found, will use system CA certificates")
		return nil
	}

	// Load client certificate
	certificate, err := tls.LoadX509KeyPair(ClientCertificatePath, ClientKeyPath)
	if err != nil {
		log.WithError(err).Errorf("Couldn't load client TLS certificate")
		return nil
	}

	log.Debugf("Creating a TLS client")
	// Create client with TLS config
	return ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{certificate},
	})
}
