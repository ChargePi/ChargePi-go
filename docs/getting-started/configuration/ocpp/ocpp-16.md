# OCPP 1.6

The default/reference configuration for protocol version 1.6 can be found [here](../../../../configs/ocpp-config.yaml).

Each OCPP 1.6 configuration variable is represented as a dictionary with key equal to **variable name**, the **value** and
**permission** attributes. For more information regarding OCPP 1.6 configuration,
visit [this link](https://www.oasis-open.org/committees/download.php/58944/ocpp-1.6.pdf).

```json
{
  "version": 1,
  "keys": [
    {
      "key": "AllowOfflineTxForUnknownId",
      "readOnly": false,
      "value": "false"
    },
    {
      "key": "AuthorizationCacheEnabled",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "AuthorizeRemoteTxRequests",
      "readOnly": false,
      "value": "false"
    },
    {
      "key": "ClockAlignedDataInterval",
      "readOnly": false,
      "value": "0"
    },
    {
      "key": "ConnectionTimeOut",
      "readOnly": false,
      "value": "50"
    },
    {
      "key": "GetConfigurationMaxKeys",
      "readOnly": false,
      "value": "30"
    },
    {
      "key": "HeartbeatInterval",
      "readOnly": false,
      "value": "60"
    },
    {
      "key": "LocalAuthorizeOffline",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "LocalPreAuthorize",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "MaxEnergyOnInvalidId",
      "readOnly": false,
      "value": "0"
    },
    {
      "key": "MeterValuesSampledData",
      "readOnly": false,
      "value": "Power.Active.Import"
    },
    {
      "key": "MeterValuesAlignedData",
      "readOnly": false,
      "value": "false"
    },
    {
      "key": "NumberOfConnectors",
      "readOnly": false,
      "value": "6"
    },
    {
      "key": "MeterValueSampleInterval",
      "readOnly": false,
      "value": "60"
    },
    {
      "key": "ResetRetries",
      "readOnly": false,
      "value": "3"
    },
    {
      "key": "ConnectorPhaseRotation",
      "readOnly": false,
      "value": "0.RST, 1.RST, 2.RTS"
    },
    {
      "key": "StopTransactionOnEVSideDisconnect",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "StopTransactionOnInvalidId",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "StopTxnAlignedData",
      "readOnly": false,
      "value": ""
    },
    {
      "key": "StopTxnSampledData",
      "readOnly": false,
      "value": ""
    },
    {
      "key": "SupportedFeatureProfiles",
      "readOnly": false,
      "value": "Core, LocalAuthListManagement, Reservation, RemoteTrigger"
    },
    {
      "key": "TransactionMessageAttempts",
      "readOnly": false,
      "value": "3"
    },
    {
      "key": "TransactionMessageRetryInterval",
      "readOnly": false,
      "value": "60"
    },
    {
      "key": "UnlockConnectorOnEVSideDisconnect",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "ReserveConnectorZeroSupported",
      "readOnly": false,
      "value": "false"
    },
    {
      "key": "SendLocalListMaxLength",
      "readOnly": false,
      "value": "20"
    },
    {
      "key": "LocalAuthListEnabled",
      "readOnly": false,
      "value": "true"
    },
    {
      "key": "LocalAuthListMaxLength",
      "readOnly": false,
      "value": "20"
    }
  ]
}
```