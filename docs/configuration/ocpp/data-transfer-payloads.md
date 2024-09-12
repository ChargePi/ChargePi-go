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

### Displaying messages

The ChargePi supports remote display management, allowing to remotely change screens on the charge point.
For forward-compatibility, all requests in OCPP 1.6 Data Transfer are the same as in OCPP 2.0.1 specification.

The VendorID must match the vendor id on the charge point. The message ID is the same as declared in the
OCPP 2.0.1 specification.