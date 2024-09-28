//go:build raspberrypi4

package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitString(t *testing.T) {
	tests := []struct {
		name          string
		combinedLines string
		size          int
		want          []string
	}{
		{
			name:          "combinedLines is empty",
			combinedLines: "",
			size:          16,
			want:          []string{""},
		},
		{
			name:          "combinedLines is less than size",
			combinedLines: "Hello",
			size:          16,
			want:          []string{"Hello"},
		},
		{
			name:          "combinedLines is equal to size",
			combinedLines: "Hello, World!",
			size:          13,
			want:          []string{"Hello, World!"},
		},
		{
			name:          "combinedLines is greater than size",
			combinedLines: "Hello, World!",
			size:          5,
			want:          []string{"Hello", ", Wor", "ld!"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := splitString(test.combinedLines, size)

			assert.Equal(t, test.want, got)
		})
	}
}
