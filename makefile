.PHONY: install-dependencies buf-lint buf-generate format lint

install-dependencies:
	sudo sh ./scripts/install-dependencies.sh pn532_uart 0

buf-lint:
	buf lint

buf-generate:
	buf generate

format:
	golangci-lint run --fix

lint: buf-lint
	golangci-lint run
