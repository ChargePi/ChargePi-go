.PHONY: install-dependencies buf-lint buf-generate format lint gen-mocks generate lint gci go-fmt format

install-dependencies:
	sudo sh ./scripts/install-dependencies.sh pn532_uart 0

generate: buf-generate gen-mocks

gen-mocks:
	mockery

buf-generate:
	buf generate

format: go-fmt buf-format

buf-format:
	buf format -w

go-fmt:
	golangci-lint run --fix

lint: buf-lint gci

gci:
	golangci-lint run

buf-lint:
	buf lint
