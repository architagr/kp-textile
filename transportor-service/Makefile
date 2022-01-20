.PHONY: vendor
vendor:
	go mod vendor

test:
	go test -race ./...

build:
	go build -mod=vendor .

lambda-build:
	rm -f internal-api main.zip
	GOOS=linux go build clients/internal-api.go
	zip main.zip internal-api
	rm -f internal-api

start-internal:
	go run ./clients/internal-api.go

