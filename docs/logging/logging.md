# 📜 Logging

Logging is configured in the `settings` file. There are multiple ways of logging:

1. Remote logging
    - Using Graylog and GELF
    - Using Syslog

2. File logging
    - Log to the `/var/logs` in a JSON format
    - Read the logs from the file and send them to a remote Loki instance using Promtail (or any other client)

3. Console logging
    - Only output the logs in the console.

## 🛠️ Logging configuration

| Attribute |        Valid values         |  Default  |                     Description                      |
|:---------:|:---------------------------:|:---------:|:----------------------------------------------------:|
|  `type`   | `remote`, `file`, `console` | `console` | Where to output the logs. Can have multiple outputs. |
| `format`  |  `json`, `syslog`, `gelf`   |  `json`   |             The format the logs are in.              |
|  `host`   |    Any valid IP/hostname    |     /     |        Only needed when the type is `remote`.        |
|  `port`   |              /              |     /     |        Only needed when the type is `remote`.        |

Example logging settings:

```json
{
  // ... other settings
  "logging": {
    "type": [
      "remote",
      "file"
    ],
    "format": "gelf",
    "host": "logging.example.com",
    "port": 12201
  }
}
```