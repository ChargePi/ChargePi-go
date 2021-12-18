package tls

import (
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

const (
	InvalidCACertificatePath = "../../test/certs/cs/ca123.crt"
	CACertificatePath        = "../../test/certs/ca.crt"
	ClientCertificatePath    = "../../test/certs/cp/charge-point.crt"
	ClientKeyPath            = "../../test/certs/cp/charge-point.key"
)

func Test_getTLSClient(t *testing.T) {
	var (
		require = require.New(t)
	)

	exec.Command("sudo ../../test/create-test-certs.sh")

	// Invalid paths
	require.Nil(GetTLSClient(InvalidCACertificatePath, ClientCertificatePath, ClientKeyPath))
	require.Nil(GetTLSClient(CACertificatePath, "certs/invalidCertificatePath.crt", ClientKeyPath))
	require.Nil(GetTLSClient(CACertificatePath, ClientCertificatePath, "certs/cp/charge-point-invalid.key"))

	// Valid combination
	require.NotNil(GetTLSClient(CACertificatePath, ClientCertificatePath, ClientKeyPath))
}
