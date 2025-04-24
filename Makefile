wire:
	wire gen ./internal/app/

lint:
	golangci-lint run --out-format=colored-line-number