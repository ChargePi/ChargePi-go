# OCPP 2.0.1

Configuration for protocol version 2.0.1 can be found in _/charge_point/v201/configuration/configuration.json_.

In the protocol version 2.0.1, configuration variables are nested in Controllers (postfix - Ctrlr). Each controller has
variables represented as a dictionary with attributes: **readOnly**, **value** and _optionally_ **unit**. Some
controllers aren't completely supported. For more information regarding OCPP 2.0.1 configuration,
visit [the official website](https://www.openchargealliance.org/protocols/ocpp-201/).

```json
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