# Supported hardware

## Note:

The hardware, such as Reader, Indicator and a Display should be configured in [_settings
file_](../../../configs/settings.json), and EVSEs should be configured in the [_evses folder_](../../../configs/evses),
each in a separate file.

## RFID/NFC readers

| Reader | Is supported | 
|:------:|:------------:|
| PN532  |      ✔       |

## Displays

| Display | Is supported | 
|:-------:|:------------:|
| HD44780 |      ✔       |

## EVCC

EV charging controller (EVCC) controls the communication with the EV and allows or denies the charging. It can also set
the charging current limit.

|         EVCC          | Is supported | 
|:---------------------:|:------------:|
|         Relay         |      ✔       |
| Phoenix Contact EVSEs |   Planned    |

## Power meters

| Power meter | Is supported | 
|:-----------:|:------------:|
|   CS5460A   |      ✔       |
|     ETI     |   Planned    |

## Indicators

| Indicator | Is supported | 
|:---------:|:------------:|
|  WS2812b  |      ✔       |
|  WS2811   |      ✔       |

## Contributing

If you want to add support for any type of hardware, read the
contribution [guide](../../contribution/hardware/adding-support-for-hardware.md).