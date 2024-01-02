# Importing and exporting configuration(s)

There are three types of configurations:

1. [charge point](./configuration.md/#-charge-point-general-configuration)
2. [ocpp](./ocpp)
3. [evse](evse-configuration.md#the-evse-configuration)

Each configuration can be imported and exported separately. Use the following commands to import and export
configurations:

## Importing

```bash
chargepi import <flag> <path-to-configuration-file>
```

Available flags:

|   Flag   |                                Description                                | Default value | 
|:--------:|:-------------------------------------------------------------------------:|:-------------:|
|  --evse  |                      EVSE configuration folder path                       |               |
|  --ocpp  | OCPP configuration file path. Requires the --version flag to also be set. |               |
| --config |                     ChargePi configuration file path                      |               |

## Exporting

```bash
chargepi import <configuration-type> <path-to-configuration-file>
```