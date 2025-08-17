package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type exampleSettings struct {
	ExampleField    string           `json:"example_field"`
	ExampleSettings exampleSettings2 `json:"example_settings_2"`
}

type exampleSettings2 struct {
	ExampleField string `json:"example_field"`
}

func TestInitSettings(t *testing.T) {
	tests := []struct {
		name         string
		panics       bool
		settingsPath string
	}{
		{
			name: "Settings file found in path",
		},
		{
			name: "Settings file path empty but default settings file found",
		},
		{
			name: "Settings path empty and default settings file not found",
		},
		{
			name: "Settings path not found",
		},
		{
			name: "Settings file unexpected format",
		},
		{
			name: "Settings file not formatted correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				assert.Panics(t, func() {

				})
				return
			}
		})
	}
}

func Test_readConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		panics       bool
		settingsPath string
	}{
		{
			name: "Settings file found in path",
		},
		{
			name: "Settings file path empty but default settings file found",
		},
		{
			name: "Settings path empty and default settings file not found",
		},
		{
			name: "Settings path not found",
		},
		{
			name: "Settings file unexpected format",
		},
		{
			name: "Settings file not formatted correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func Test_setupEnv(t *testing.T) {
	tests := []struct {
		name         string
		panics       bool
		settingsPath string
	}{
		{
			name: "Settings file found in path",
		},
		{
			name: "Settings file path empty but default settings file found",
		},
		{
			name: "Settings path empty and default settings file not found",
		},
		{
			name: "Settings path not found",
		},
		{
			name: "Settings file unexpected format",
		},
		{
			name: "Settings file not formatted correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func Test_setDefaults(t *testing.T) {
	tests := []struct {
		name         string
		panics       bool
		settingsPath string
	}{
		{
			name: "Settings file found in path",
		},
		{
			name: "Settings file path empty but default settings file found",
		},
		{
			name: "Settings path empty and default settings file not found",
		},
		{
			name: "Settings path not found",
		},
		{
			name: "Settings file unexpected format",
		},
		{
			name: "Settings file not formatted correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
