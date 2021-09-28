package reader

import (
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
)

const (
	PN532 = "PN532"
)

type Reader interface {
	init()
	ListenForTags()
	Cleanup()
	Reset()
	GetTagChannel() chan string
}

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
