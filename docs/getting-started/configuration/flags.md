# 🚩 Flags and environment variables

Here is a list of all available flags:

|        Flag         | Short |           Description           | Default value |
|:-------------------:|:-----:|:-------------------------------:|:-------------:|
|     `-settings`     |   /   |   Path to the settings file.    |               |
| `-connector-folder` |   /   |  Path to the connector folder.  |               |
|   `-ocpp-config`    |   /   | Path to the OCPP configuration. |               |
|       `-auth`       |   /   | Path to the authorization file. |               |
|      `-debug`       | `--d` |           Debug mode            |     false     |
|       `-api`        | `--a` |         Expose the API          |     false     |
|   `-api-address`    |   /   |           API address           |  "localhost"  |
|     `-api-port`     |   /   |            API port             |     4269      |

Environment variables are created automatically thanks to [Viper](https://github.com/spf13/viper) and are prefixed
with `CHARGEPI`. Only the settings file attributes are bound to the env variables as well as debug mode and API
settings. Connectors do not have their attributes bound to environment variables.

Example environment variable: `CHARGEPI_CHARGEPOINT_CONNECTIONSETTINGS_ID`.
