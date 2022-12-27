//go:build !linux
// +build !linux

package reader

import (
	"context"
)

type TagReader struct {
	TagChannel    chan string
	DeviceAddress string
	ResetPin      int
}

// init Initialize the NFC/RFID tag reader. Establish the connection and set up the reader.
func (reader *TagReader) init() {
}

// ListenForTags poll the reader for NFC/RFID tags. Uses multiple modulations for different standards.
// Send the ID of the detected card through the TagChannel. If there is a problem with the reader,
// hardware Reset the device.
func (reader *TagReader) ListenForTags(ctx context.Context) {
}

func (reader *TagReader) GetTagChannel() <-chan string {
	return reader.TagChannel
}

// Cleanup Close the reader device connection.
func (reader *TagReader) Cleanup() {
	close(reader.TagChannel)
}

// Reset Implements the hardware reset by pulling the ResetPin low and then releasing.
func (reader *TagReader) Reset() {
}
