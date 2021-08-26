package hardware

import (
	"encoding/hex"
	"fmt"
	"github.com/clausecker/nfc/v2"
	"github.com/warthog618/gpiod"
	"log"
	"time"
)

type TagReader struct {
	TagChannel       chan string
	reader           *nfc.Device
	DeviceConnection string
	ResetPin         int
}

// init Initialize the NFC/RFID tag reader. Establish UART connection and setup the reader.
func (reader *TagReader) init() {
	dev, err := nfc.Open(reader.DeviceConnection)
	if err != nil {
		log.Fatalf("Cannot communicate with NFC reader")
		return
	}
	reader.reader = &dev
	err = reader.reader.InitiatorInit()
	if err != nil {
		log.Fatal("Failed to initialize")
		return
	}
}

// ListenForTags poll the reader for NFC/RFID tags. Uses multiple modulations for different standards.
// Send the ID of the detected card through the TagChannel. If there is a problem with the reader,
// hardware Reset the device.
func (reader *TagReader) ListenForTags() {
	reader.init()
	var modulations = []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106},
		{Type: nfc.Felica, BaudRate: nfc.Nbr212},
		{Type: nfc.Felica, BaudRate: nfc.Nbr424},
		{Type: nfc.Jewel, BaudRate: nfc.Nbr106},
		{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106},
	}
	var (
		err    error
		count  int
		target nfc.Target
		UID    string
	)
	for {
		count, target, err = reader.reader.InitiatorPollTarget(modulations, 1, 300*time.Millisecond)
		if err != nil {
			fmt.Println("Error polling the reader", err)
			reader.Reset()
			continue
		}
		if count > 0 {
			fmt.Println(target.String())
			switch target.Modulation() {
			case nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443aTarget)
				var UIDLen = card.UIDLen
				var ID = card.UID
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.ISO14443b, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443bTarget)
				var UIDLen = len(card.ApplicationData)
				var ID = card.ApplicationData
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr212}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Felica, BaudRate: nfc.Nbr424}:
				var card = target.(*nfc.FelicaTarget)
				var UIDLen = card.Len
				var ID = card.ID
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.Jewel, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.JewelTarget)
				var ID = card.ID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			case nfc.Modulation{Type: nfc.ISO14443biClass, BaudRate: nfc.Nbr106}:
				var card = target.(*nfc.ISO14443biClassTarget)
				var ID = card.UID
				var UIDLen = len(ID)
				UID = hex.EncodeToString(ID[:UIDLen])
				break
			}
			reader.TagChannel <- UID
		}
		time.Sleep(time.Second)
	}
}

// Cleanup Close the reader device connection.
func (reader *TagReader) Cleanup() {
	reader.reader.Close()
}

// Reset Implements the hardware reset by pulling the ResetPin low and then releasing.
func (reader *TagReader) Reset() {
	log.Println("Resetting the reader..")
	//refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	pin, err := c.RequestLine(reader.ResetPin, gpiod.AsOutput(0))
	if err != nil {
		log.Println(err)
		return
	}
	err = pin.SetValue(0)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 100)
	err = pin.SetValue(1)
	if err != nil {
		log.Println(err)
		return
	}
}
