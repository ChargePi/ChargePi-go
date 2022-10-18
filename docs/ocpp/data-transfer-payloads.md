# Custom data payloads

## OCPP 1.6

### Charge point model additional info

This payload contains the type of the charge point and it's maximum power output for all EVSEs.
It is sent to the Central System after it is approved by the Central System.

```json
{
  "type": "AC",
  "maxPower": 6.0
}
```

### Connector details

This payload contains maximum power output for an EVSE and the types of connectors the EVSE has.
It is sent to the Central System after loading all the EVSE settings to the Charge point.

```json
{
  "evseId": 1,
  "maxPower": 6.0,
  "connectors": [
    {
      "id": 1,
      "type": "CCS-1"
    },
    {
      "id": 2,
      "type": "ChaDeMo"
    }
  ]
}
```