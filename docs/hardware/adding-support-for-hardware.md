# Adding support for hardware

There are four hardware component groups that are included in the project:

1. NFC/RFID tag reader,
2. LCD (display),
3. (Led) Indicator,
4. Power meter

These hardware components have corresponding interfaces that are included in the `ChargePointHandler` struct. This
allows adding support for other models of hardware with similar functionalities.

## Reader hardware

All readers must implement the `Reader` interface. It is recommended that you implement the interface in a new file
named after the model of the reader in the `hardware/reader` package. Then you should add a **constant** named after
the **model** of the reader in the `reader` file in the package and add a switch case with the implementation and the
necessary logic that returns a pointer to the struct.

The settings of the reader are read from the `settings.json` file, which is stored in the cache and are available in the
NewTagReader method.

```golang
package reader

const (
	// Add the reader model here
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
	//...
	if tagReaderSettings.IsSupported {
		log.Println("Preparing tag reader from config:", tagReaderSettings.ReaderModel)
		switch tagReaderSettings.ReaderModel {
		// Add a new case with your implementation and return the pointer
		case PN532:
			//...
			return reader
		default:
			return nil
		}
	}
	return nil
}
```

## Display hardware

All displays must implement the `LCD` interface. It is recommended that you implement the interface in a new file named
after the model of the display/LCD in the `hardware/display` package. Then you should add a **constant** named after
the **model** of the display in the `display` file in the package and add a switch case with the implementation and the
necessary logic that returns a pointer to the struct.

```golang
package display

const (
	// Add the LCD driver here
	DriverHD44780 = "hd44780"
)

type (
	// LCDMessage Object representing the message that will be displayed on the LCD.
	// Each array element in Messages represents a line being displayed on the 16x2 screen.
	LCDMessage struct {
		Messages        []string
		messageDuration int
	}

	// LCD is an abstraction layer for concrete implementation of a display.
	LCD interface {
		DisplayMessage(message LCDMessage)
		ListenForMessages()
		Cleanup()
		Clear()
		GetLcdChannel() chan LCDMessage
	}
)

// NewDisplay returns a concrete implementation of an LCD based on the drivers that are supported.
// The LCD is built with the settings from the settings file.
func NewDisplay() LCD {
	//...
	if lcdSettings.IsSupported {
		log.Println("Preparing LCD from config")
		switch lcdSettings.Driver {
		// Add a new case with your implementation and return the pointer
		case DriverHD44780:
			//..
			return lcd
		default:
			return nil
		}
	}
	return nil
}
```

## Indicator hardware

The process is the same as the previous description.

```golang
package indicator

const (
	//...
	// Add your indicator const here
	TypeWS281x = "WS281x"
)

type Indicator interface {
	DisplayColor(index int, colorHex uint32) error
	Blink(index int, times int, colorHex uint32) error
	Cleanup()
}

// NewIndicator constructs the Indicator based on the type provided by the settings file.
func NewIndicator(stripLength int) Indicator {
	//...
	if indicatorSettings.Enabled {
		log.Println("Preparing Indicator from config: ", indicatorSettings.Type)
		switch indicatorSettings.Type {
		// Add a case with your implementation here
		case TypeWS281x:
			//...
			return ledStrip
		default:
			return nil
		}
	}
	return nil
}
```

## Power meters

The process is the same as the previous description.

```golang
package power_meter

const (
	// Add your power meter type here
	TypeC5460A = "cs5460a"
)

type PowerMeter interface {
	Reset()
	GetEnergy() float64
	GetPower() float64
	GetCurrent() float64
	GetVoltage() float64
	GetRMSCurrent() float64
	GetRMSVoltage() float64
}

func NewPowerMeter(connector *settings.Connector) (PowerMeter, error) {
	if connector.PowerMeter.Enabled {
		log.Println("Creating a new power meter:", connector.PowerMeter.Type)
		switch connector.PowerMeter.Type {
		// Add your case with implementation here
		case TypeC5460A:
			return NewCS5460PowerMeter(
				connector.PowerMeter.PowerMeterPin,
				connector.PowerMeter.SpiBus,
				connector.PowerMeter.ShuntOffset,
				connector.PowerMeter.VoltageDividerOffset,
			)
		default:
			return nil, fmt.Errorf("power meter type not supported")
		}
	}
	return nil, fmt.Errorf("power meter not enabled")
}
```

