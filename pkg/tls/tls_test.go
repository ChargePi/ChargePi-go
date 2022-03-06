package tls

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

const (
	CACertificatePath     = "./certs/ca.crt"
	ClientCertificatePath = "./certs/cp/charge-point.crt"
	ClientKeyPath         = "./certs/cp/charge-point.key"

	InvalidCACertificatePath     = "./certs/cs/ca123.crt"
	InvalidClientCertificatePath = "./certs/invalidCertificatePath.crt"
	InvalidClientKeyPath         = "./certs/cp/charge-point-invalid.key"
)

var (
	wd, _      = os.Getwd()
	scriptPath = wd + "/../../test/tools"
	scriptName = "create-test-certs.sh"
)

func Test_getTLSClient(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert2.New(t)

		script = fmt.Sprintf("%s/%s", scriptPath, scriptName)
		cmd    = exec.Command("/bin/sh", script, scriptPath)
	)

	err := cmd.Run()
	require.NoError(err)

	// Invalid paths
	assert.Nil(GetTLSClient(InvalidCACertificatePath, ClientCertificatePath, ClientKeyPath))
	assert.Nil(GetTLSClient(CACertificatePath, InvalidClientCertificatePath, ClientKeyPath))
	assert.Nil(GetTLSClient(CACertificatePath, ClientCertificatePath, InvalidClientKeyPath))

	// Valid combination
	assert.NotNil(GetTLSClient(CACertificatePath, ClientCertificatePath, ClientKeyPath))
}
