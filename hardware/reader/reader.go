package reader

import (
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
)

// supported readers
const (
	PN532 = "PN532"
)

// Reader is an abstraction for an RFID/NFC tag reader.
type Reader interface {
	ListenForTags()
	Cleanup()
	Reset()
	GetTagChannel() chan string
}

// NewTagReader creates an instance of the Reader interface based on the provided configuration.
func NewTagReader() Reader {
	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config := cacheSettings.(*settings.Settings)
	tagReaderSettings := config.ChargePoint.Hardware.TagReader

	if tagReaderSettings.IsSupported {
		log.Println("Preparing tag reader from config:", tagReaderSettings.ReaderModel)
		switch tagReaderSettings.ReaderModel {
		case PN532:
			tagChannel := make(chan string)
			tagReader := &TagReader{
				TagChannel:       tagChannel,
				DeviceConnection: tagReaderSettings.Device,
				ResetPin:         tagReaderSettings.ResetPin,
			}
			return tagReader
		default:
			return nil
		}
	}
	return nil
}
