# Client API

The client comes with a GRPC API that enables other services to integrate with the client. One such use-case would be
writing a frontend application that interacts with the client and displays the values for an ongoing transaction or
starts/stops a transaction.

To access the endpoint, it must be enabled through flags (`--a`) and will be exposed by default
on `localhost:4269`.

## Endpoints

Compiling the protobuf:

```bash
protoc --go_out=./internal/api --go_opt=paths=source_relative \
    --proto_path=pkg/proto \
    --go-grpc_out=./internal/api --go-grpc_opt=paths=source_relative \
    pkg/proto/*.proto
```
