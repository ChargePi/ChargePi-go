## Configuration variables

### For OCPP 1.6

Configuration for protocol version 1.6 can be found in _/configs/configuration.json_.

Each OCPP 1.6 configuration variable is represented as a dictionary with key equal to variable name, the **value** and
**permission** attributes. For more information regarding OCPP 1.6 configuration,
read the [official documentation](https://www.oasis-open.org/committees/download.php/58944/ocpp-1.6.pdf).

```
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
    },
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
    },
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

### For OCPP 2.0.1

Configuration for protocol version 2.0.1 can be found in _/configs/configuration.json_.

In the protocol version 2.0.1, configuration variables are nested in Controllers (postfix - Ctrlr). Each controller has
variables represented as a dictionary with attributes: **readOnly**, **value** and _optionally_ **unit**. Some
controllers aren't completely supported. For more information regarding OCPP 2.0.1 configuration,
visit [the official website](https://www.openchargealliance.org/protocols/ocpp-201/).

```
{
  "AlignedDataCtrlr": {
    "Enabled": false,
    "Measurands": {
      "readOnly": false,
      "value": []
    },
    "Interval": {
      "readOnly": false,
      "value": 50,
      "unit": "seconds"
    },
    "TxEndedMeasurands": {
      "readOnly": false,
      "value": []
    },
    "TxEndedInterval": {
      "readOnly": false,
      "value": 50,
      "unit": "seconds"
    }
  },
  "AuthCacheCtrlr": {
    "Enabled": false,
    "AuthCacheEnabled": {
      "readOnly": false,
      "value": true
    },
    "AuthCacheLifeTime": {
      "readOnly": false,
      "value": true
    }
  },
  "AuthCtrlr": {
    "Enabled": true,
    "AuthorizeRemoteStart": {
      "readOnly": true,
      "value": true
    },
    "AuthEnabled": {
      "readOnly": false,
      "value": true
    },
    "OfflineTxForUnknownIdEnabled": {
      "readOnly": false,
      "value": false
    },
    "LocalAuthorizeOffline": {
      "readOnly": false,
      "value": true
    },
    "LocalPreAuthorize": {
      "readOnly": false,
      "value": true
    }
  },
  "CHAdeMOCtrlr": {
    "Enabled": false
  },
  "ClockCtrlr": {
    "Enabled": false,
    "DateTime": {
      "readOnly": true,
      "value": 1
    },
    "TimeSource": {
      "readOnly": false,
      "value": [
        "Heartbeat",
        "MobileNetwork"
      ]
    }
  },
  "CustomizationCtrlr": {
    "Enabled": false
  },
  "DeviceDataCtrlr": {
    "Enabled": false,
    "ItemsPerMessage": {
      "readOnly": true,
      "value": 1
    },
    "BytesPerMessage": {
      "readOnly": true,
      "value": 1
    }
  },
  "DisplayMessageCtrlr": {
    "Enabled": false
  },
  "ISO15118Ctrlr": {
    "Enabled": false
  },
  "LocalAuthListCtrlr": {
    "Enabled": false
  },
  "MonitoringCtrlr": {
    "Enabled": false
  },
  "OCPPCommCtrlr": {
    "Enabled": true,
    "RetryBackOffRepeatTimes": {
      "readOnly": false,
      "value": 1
    },
    "RetryBackOffRandomRange": {
      "readOnly": false,
      "unit": "seconds",
      "value": 1
    },
    "RetryBackOffWaitMinimum": {
      "readOnly": false,
      "unit": "seconds",
      "value": 1
    },
    "WebSocketPingInterval": {
      "readOnly": false,
      "unit": "seconds",
      "value": 1
    },
    "DefaultMessageTimeout": {
      "readOnly": true,
      "unit": "seconds",
      "value": 1
    },
    "FileTransferProtocols": {
      "readOnly": true,
      "value": [
        "HTTP",
        "HTTPS"
      ]
    },
    "HeartbeatInterval": {
      "readOnly": false,
      "unit": "seconds",
      "value": 60
    },
    "NetworkConfigurationPriority": {
      "readOnly": false,
      "value": []
    },
    "NetworkProfileConnectionAttempts": {
      "readOnly": false,
      "value": 1
    },
    "OfflineThreshold": {
      "readOnly": false,
      "unit": "seconds",
      "value": 150
    },
    "MessageAttempts": {
      "readOnly": false,
      "value": 1
    },
    "MessageAttemptInterval": {
      "unit": "seconds",
      "value": 90,
      "readOnly": false
    },
    "UnlockOnEVSideDisconnect": {
      "value": true,
      "readOnly": false
    },
    "ResetRetries": {
      "readOnly": false,
      "value": 90
    }
  },
  "ReservationCtrlr": {
    "Enabled": false
  },
  "SampledDataCtrlr": {
    "Enabled": true,
    "SampledDataEnabled": {
      "readOnly": false,
      "value": true
    },
    "TxEndedMeasurands": {
      "readOnly": false,
      "value": []
    },
    "TxEndedInterval": {
      "readOnly": false,
      "value": 60,
      "unit": "seconds"
    },
    "TxStartedMeasurands": {
      "readOnly": false,
      "value": []
    },
    "TxUpdatedMeasurands": {
      "readOnly": false,
      "value": []
    },
    "TxUpdatedInterval": {
      "readOnly": false,
      "value": []
    }
  },
  "SecurityCtrlr": {
    "Enabled": false,
    "OrganizationName": {
      "readOnly": false,
      "value": "UL FE"
    },
    "CertificateEntries": {
      "readOnly": true,
      "value": 1
    },
    "SecurityProfile": {
      "readOnly": true,
      "value": 1
    }
  },
  "SmartChargingCtrlr": {
    "Enabled": false
  },
  "TariffCostCtrlr": {
    "Enabled": false
  },
  "TxCtrlr": {
    "Enabled": true,
    "EVConnectionTimeOut": {
      "readOnly": false,
      "value": 60,
      "unit": "seconds"
    },
    "StopTxOnEVSideDisconnect": {
      "readOnly": true,
      "value": true
    },
    "TxStartPoint": {
      "readOnly": true,
      "value": [
        "Authorized"
      ]
    },
    "TxStopPoint": {
      "readOnly": true,
      "value": [
        "PowerPathClosed"
      ]
    },
    "StopTxOnInvalidId": {
      "readOnly": false,
      "value": true
    }
  }
}
```
