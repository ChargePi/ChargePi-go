//go:build linux
// +build linux

package reader

import (
	"context"
	"encoding/hex"
	"go.uber.org/zap"
	"time"

	"github.com/clausecker/nfc/v2"
	"github.com/warthog618/gpiod"
)

var (
	Iso14443A  = nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}
	ISO14443b  = nfc.Modulation{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106}
	Felica1    = nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr212}
	Felica2    = nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr424}
	Jewel      = nfc.Modulation{Type: nfc.Jewel, BaudRate: nfc.Nbr106}
	ISO14443bi = nfc.Modulation{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106}

	modulations = []nfc.Modulation{
		Iso14443A, ISO14443b, Felica1, Felica2, Jewel, ISO14443bi,
	}
)

type TagReader struct {
	tagChannel chan string
	reader     *nfc.Device
	devAddress string
	resetPin   int
	deviceType string
	logger     *zap.Logger
}

func NewReader(logger *zap.Logger, device, deviceType string, resetPin int) (*TagReader, error) {
	tagChannel := make(chan string, 5)

	return &TagReader{
		logger:     logger.Named("libnfc_tag_reader"),
		tagChannel: tagChannel,
		devAddress: device,
		resetPin:   resetPin,
		deviceType: deviceType,
	}, nil
}

// init Initialize the NFC/RFID tag reader. Establish the connection and set up the reader.
func (reader *TagReader) init() {
	dev, err := nfc.Open(reader.devAddress)
	if err != nil {
		reader.logger.Panic("Cannot communicate with reader")
	}

	reader.reader = &dev
	err = reader.reader.InitiatorInit()
	if err != nil {
		reader.logger.Panic("Failed to initialize reader")
	}
}

// ListenForTags poll the reader for NFC/RFID tags. Uses multiple modulations for different standards.
// Send the ID of the detected card through the TagChannel. If there is a problem with the reader,
// hardware Reset the device.
func (reader *TagReader) ListenForTags(ctx context.Context) {
	reader.init()
	var (
		err    error
		count  int
		target nfc.Target
		UID    string
		UIDLen int
	)

Listener:
	for {
		select {
		case <-ctx.Done():
			break Listener
		default:
			count, target, err = reader.reader.InitiatorPollTarget(modulations, 1, 300*time.Millisecond)
			if err != nil {
				reader.logger.With(zap.Error(err)).Error("Error polling the reader")
				reader.Reset()
				continue
			}

			if count > 0 {
				switch target.Modulation() {
				case Iso14443A:
					var card = target.(*nfc.ISO14443aTarget)
					UIDLen = card.UIDLen
					var ID = card.UID
					UID = hex.EncodeToString(ID[:UIDLen])

				case ISO14443b:
					var card = target.(*nfc.ISO14443bTarget)
					UIDLen = len(card.ApplicationData)
					var ID = card.ApplicationData
					UID = hex.EncodeToString(ID[:UIDLen])

				case Felica1, Felica2:
					var card = target.(*nfc.FelicaTarget)
					var UIDLen = card.Len
					var ID = card.ID
					UID = hex.EncodeToString(ID[:UIDLen])

				case Jewel:
					var card = target.(*nfc.JewelTarget)
					var ID = card.ID
					UIDLen = len(ID)
					UID = hex.EncodeToString(ID[:UIDLen])

				case ISO14443bi:
					var card = target.(*nfc.ISO14443biClassTarget)
					var ID = card.UID
					UIDLen = len(ID)
					UID = hex.EncodeToString(ID[:UIDLen])
				default:
					reader.logger.Warn("NFC modulation unknown")
					continue
				}

				reader.tagChannel <- UID
			}

			time.Sleep(time.Millisecond * 200)
		}
	}
}

func (reader *TagReader) GetTagChannel() <-chan string {
	return reader.tagChannel
}

func (reader *TagReader) GetType() string {
	return reader.deviceType
}

// Cleanup Close the reader device connection.
func (reader *TagReader) Cleanup() {
	close(reader.tagChannel)
	err := reader.reader.Close()
	if err != nil {
		reader.logger.With(zap.Error(err)).Error("Error closing the reader")
	}
}

// Reset Implements the hardware reset by pulling the resetPin low and then releasing.
func (reader *TagReader) Reset() {
	reader.logger.Info("Resetting the reader")

	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	pin, err := c.RequestLine(reader.resetPin, gpiod.AsOutput(1))
	if err != nil {
		reader.logger.With(zap.Error(err)).Error("Error requesting the reset line")
		return
	}

	time.Sleep(time.Millisecond * 300)
	err = pin.SetValue(0)
	if err != nil {
		reader.logger.With(zap.Error(err)).Error("Error requesting the reset line")
		return
	}

	time.Sleep(time.Millisecond * 300)
	err = pin.SetValue(1)
	if err != nil {
		reader.logger.With(zap.Error(err)).Error("Error requesting the reset line")
		return
	}
}
