package settings

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	require := assert.New(t)

	Test123 := struct {
		Enabled bool   `json:"enabled"`
		Type    string `json:"type"`
	}{
		Enabled: false,
		Type:    "",
	}

	err := WriteToFile("test123.json", &Test123)
	require.NoError(err)
	require.FileExists("test123.json")

	err = WriteToFile("test123.yaml", &Test123)
	require.NoError(err)
	require.FileExists("test123.yaml")

	err = WriteToFile("test123.o", &Test123)
	require.Error(err)
	require.NoFileExists("test123.o")

	// Clean up
	cmd := exec.Command("rm", "test123.*")
	err = cmd.Run()
}
