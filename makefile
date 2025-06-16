.PHONY: install-dependencies buf-lint buf-generate

install-dependencies:
	sudo sh ./scripts/install-dependencies.sh pn532_uart 0

buf-lint:
	buf lint

buf-generate:
	buf generate