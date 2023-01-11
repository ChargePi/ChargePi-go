proto:
	protoc --go_out=./pkg/grpc --go_opt=paths=source_relative \
		--proto_path=pkg/proto \
		--go-grpc_out=./pkg/grpc --go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto

