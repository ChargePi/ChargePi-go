package settings

import (
	"github.com/stretchr/testify/assert"
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
}
