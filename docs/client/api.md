# Client API

The client comes with a GRPC API that enables other services to integrate with the client. One such use-case would be
writing a frontend application that interacts with the client and displays the values for an ongoing transaction or
starts/stops a transaction.

To access the endpoint, it must be enabled through flags (`--a`) and will be exposed by default
on `localhost:4269`.

## Endpoints

```protobuf
syntax = "proto3";

package api;

service ChargePoint {

  rpc GetConnectorStatus (stream GetConnectorStatusRequest) returns (stream GetConnectorStatusResponse) {}

  rpc StartTransaction (StartTransactionRequest) returns (StartTransactionResponse) {}

  rpc StopTransaction (StopTransactionRequest) returns (StopTransactionResponse) {}

  rpc HandleCharging (HandleChargingRequest) returns (HandleChargingResponse) {}
}

enum ConnectorStatus {
  Available = 0;
  Preparing = 1;
  Charging = 2;
  Finishing = 3;
  Unavailable = 4;
  SuspendedEVSE = 5;
  SuspendedEV = 6;
  Reserved = 7;
  Faulted = 8;
};

enum ErrorCode {
  NoError = 0;
  OtherError = 1;
  ConnectorLockFailure = 3;
  EVCommunicationError = 4;
  GroundFailure = 5;
  HighTemperature = 6;
  InternalError = 7;
  LocalListConflict = 8;
  OverCurrentFailure = 9;
  OverVoltage = 10;
  PowerMeterFailure = 11;
  PowerSwitchFailure = 12;
  ReaderFailure = 13;
  ResetFailure = 14;
  UnderVoltage = 15;
  WeakSignal = 16;
}

enum ConnectorType{
  TYPE1 = 0;
  TYPE2 = 1;
  SCHUKO = 2;
  CHADEMO = 3;
};
/*------------------GetConnectorStatus ------------------------ */

message GetConnectorStatusRequest {
  int32 connectorId = 1;
  int32 evseId = 2;
}

message GetConnectorStatusResponse {
  ConnectorType connectorType = 1;
  ConnectorStatus connectorStatus = 2;
  ErrorCode errorCode = 3;

  string transactionId = 4;
  int32 timeElapsed = 5;
  float energyConsumed = 6;
  float currentPower = 7;
}

/*------------------ StartTransaction ------------------------ */

message StartTransactionRequest {
  string tagId = 1;
  int32 connectorId = 2;
}

message StartTransactionResponse {
  ConnectorStatus status = 1;
  string errorMessage = 2;
  int32 connectorId = 3;
}

/*------------------ StopTransaction ------------------------ */

message StopTransactionRequest {
  string tagId = 1;
  int32 connectorId = 2;
}

message StopTransactionResponse {
  ConnectorStatus status = 1;
  string errorMessage = 2;
}

/*------------------ HandleCharging ------------------------ */

message HandleChargingRequest {
  string tagId = 1;
}

message HandleChargingResponse {
  ConnectorStatus status = 1;
  string errorMessage = 2;
  int32 connectorId = 3;
}
```