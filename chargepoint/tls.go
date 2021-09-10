package chargepoint

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/lorenzodonini/ocpp-go/ws"
	"io/ioutil"
	"log"
)

func GetTLSClient(CACertificatePath string, ClientCertificatePath string, ClientKeyPath string) *ws.Client {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caCert, err := ioutil.ReadFile(CACertificatePath)
	if err != nil {
		log.Println(err)
		return nil
	} else if !certPool.AppendCertsFromPEM(caCert) {
		log.Println("no ca.cert file found, will use system CA certificates")
		return nil
	}
	// Load client certificate
	certificate, err := tls.LoadX509KeyPair(ClientCertificatePath, ClientKeyPath)
	if err != nil {
		log.Printf("couldn't load client TLS certificate: %v \n", err)
		return nil
	}
	// Create client with TLS config
	return ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{certificate},
	})
}
