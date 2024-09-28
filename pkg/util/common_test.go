package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:all
type exampleInterface interface {
	ExampleMethod()
}

type exampleStructImplementsInterface struct {
}

func (e exampleStructImplementsInterface) ExampleMethod() {
}

func TestIsNilInterfaceOrPointer(t *testing.T) {
	tests := []struct {
		name  string
		thing interface{}
		want  bool
	}{
		{
			name:  "nil interface",
			thing: exampleInterface(nil),
			want:  true,
		},
		{
			name:  "nil pointer",
			thing: nil,
			want:  true,
		},
		{
			name:  "valid interface",
			thing: &exampleStructImplementsInterface{},
			want:  false,
		},
		{
			name:  "valid pointer",
			thing: new(int),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsNilInterfaceOrPointer(tt.thing))
		})
	}
}

func TestIsRunningInDocker(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Running in docker",
			want: true,
		},
		{
			name: "Running on host",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Running in docker" {
				_, err := os.Create("/.dockerenv")
				require.NoError(t, err)
				defer os.Remove("/.dockerenv")
			}

			assert.Equal(t, tt.want, IsRunningInContainer())
		})
	}
}
