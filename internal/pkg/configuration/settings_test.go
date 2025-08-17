package configuration

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChargePi/ChargePi-go/internal/api/grpc"
	"github.com/ChargePi/ChargePi-go/internal/api/http"
)

func TestGetRuntimeSettings(t *testing.T) {
	tests := []struct {
		name             string
		expectedSettings *RuntimeSettings
		panics           bool
	}{
		{
			name: "Settings configured correctly",
			expectedSettings: &RuntimeSettings{
				GRPC: grpc.Configuration{
					Address: "localhost:50051",
				},
				HTTP: http.Configuration{
					Address: "localhost:8080",
				},
			},
		},
		{
			name:             "Settings not configured",
			expectedSettings: nil,
			panics:           true,
		},
		{
			name:             "Settings fail validation",
			expectedSettings: nil,
			panics:           true,
		},
		{
			name:             "Settings not the same struct",
			expectedSettings: nil,
			panics:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Settings configured correctly":
				marshal, err := json.Marshal(tt.expectedSettings)
				require.NoError(t, err)

				err = viper.ReadConfig(bytes.NewBuffer(marshal))
				require.NoError(t, err)
			case "Settings not configured":
				// Skip
			case "Settings fail validation":

				marshal, err := json.Marshal(tt.expectedSettings)
				require.NoError(t, err)

				err = viper.ReadConfig(bytes.NewBuffer(marshal))
				require.NoError(t, err)
			case "Settings not the same struct":
				marshal, err := json.Marshal(http.Configuration{})
				require.NoError(t, err)

				err = viper.ReadConfig(bytes.NewBuffer(marshal))
				require.NoError(t, err)
			}

			if tt.panics {
				assert.Panics(t, func() {
					_ = GetRuntimeSettings()
				})
			} else {
				settings := GetRuntimeSettings()
				assert.Equal(t, tt.expectedSettings, settings)
			}
		})
	}
}
