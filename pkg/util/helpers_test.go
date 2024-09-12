package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	// Check if pointer is nil
	pointer := new(int)
	assert.False(t, IsNilInterfaceOrPointer(pointer))
	assert.True(t, IsNilInterfaceOrPointer(nil))

	// Check if interface is nil
	assert.False(t, IsNilInterfaceOrPointer(&exampleStructImplementsInterface{}))
}
