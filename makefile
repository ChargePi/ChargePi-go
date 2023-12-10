proto:
	mkdir -p pkg/grpc
	protoc --go_out=./pkg/grpc --go_opt=paths=source_relative \
		--proto_path=pkg/proto \
		--go-grpc_out=./pkg/grpc --go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto

install-dependencies:
	sudo sh ./scripts/install-dependencies.sh pn532_uart 0