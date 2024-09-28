package chargepoint

import (
	"context"
	"testing"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/reader"
	"github.com/stretchr/testify/assert"
)

// Wont test the logger as it is a simple wrapper around logrus
func TestWithLogger(t *testing.T) {
	t.Skipf("Skipping test for WithLogger")
}

func TestWithDisplayFromSettings(t *testing.T) {
	mock := &MockChargePoint{}

	// Display disabled
	err := mock.ApplyOpts(WithDisplayFromSettings(display.Settings{
		IsEnabled: false,
	}))
	assert.Equal(t, display.ErrDisplayDisabled, err)

	// Invalid driver
	err = mock.ApplyOpts(WithDisplayFromSettings(display.Settings{
		IsEnabled: true,
		Driver:    "invalid",
	}))
	assert.Equal(t, display.ErrDisplayUnsupported, err)

	// OK case
	err = mock.ApplyOpts(WithDisplayFromSettings(display.Settings{
		IsEnabled:    true,
		Driver:       display.TypeDummy,
		DisplayDummy: &display.DummySettings{},
	}))
	assert.Equal(t, display.ErrDisplayUnsupported, err)
}

func TestWithReaderFromSettings(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock := &MockChargePoint{}

	// Reader disabled
	err := mock.ApplyOpts(WithReaderFromSettings(ctx, reader.Settings{
		IsEnabled: false,
	}))
	assert.Equal(t, display.ErrDisplayDisabled, err)

	// Invalid reader type
	err = mock.ApplyOpts(WithReaderFromSettings(ctx, reader.Settings{
		IsEnabled:   true,
		ReaderModel: "invalid",
	}))
	assert.Equal(t, display.ErrDisplayUnsupported, err)

	// OK case
	err = mock.ApplyOpts(WithReaderFromSettings(ctx, reader.Settings{
		IsEnabled:   true,
		ReaderModel: reader.TypeDummy,
		DummyReader: &reader.DummySettings{
			TagIds: []string{"tag-id"},
		},
	}))
	assert.Equal(t, display.ErrDisplayUnsupported, err)
}

func TestWithIndicator(t *testing.T) {

}
