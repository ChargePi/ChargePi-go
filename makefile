proto:
	protoc --go_out=./internal/api/grpc --go_opt=paths=source_relative \
		--proto_path=pkg/proto \
		--go-grpc_out=./internal/api/grpc --go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto

