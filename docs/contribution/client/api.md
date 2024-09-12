# Charge point APIs

The charge point comes with a GRPC API that enables other services to integrate with the client.
One such use-case would be writing a custom frontend application that interacts with the client and
displays the values for an ongoing transaction or starts/stops a transaction.

To access the endpoint, it must be enabled through flags (`--a`) and will be exposed by default
on `localhost:4269`.

## GRPC API

Compiling the protobuf using make:

```bash
make proto
```

You can issue gRPC requests using the `grpcurl` tool:

```bash
grpcurl -plaintext -d '{"id": "1"}' localhost:4269/chargepoint.ChargePoint/GetChargePoint
```

## HTTP API

The client also exposes an HTTP API used for health checks and serving the frontend application, exposed on
`localhost:8080`.

The healthcheck is available at `/healthz` and the frontend application is available at `/`.