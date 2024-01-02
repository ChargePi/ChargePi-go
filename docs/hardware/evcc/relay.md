### Relay (or relay module)

## Description

There are multiple simple relay options that could be used for the charging station, such as Solid state relays,
contactors, etc. Choose your option according to your needs.

Be very careful with the options and consult a professional, as it may become electrical and fire hazard.

It is highly recommended splitting both GND and VCC between relays or using a relay module.

## Wiring

| RPI PIN                           | RELAY PIN | 
|-----------------------------------|:---------:|
| 4 or any 5V pin                   |    VCC    | 
| 20 or any ground pin              |    GND    |  
| 37 (GPIO 26) or any free GPIO pin | S/Enable  |  

## Example configuration

```yaml
evseId: 1
maxPower: 1.7
connectors:
  - connectorId: 1
    type: Schuko
    status: Available
evcc:
  type: Relay
  relayPin: 26
  inverseLogic: false
powerMeter:
  enabled: false
  type: CS5460A
  spi:
    bus: 0
    pin: 25
  cs5460a:
    shuntOffset: 0.055
    voltageDividerOffset: 1333
```