.PHONY: vendor
vendor:
	go mod vendor

test:
	go test -race ./...

build:
	go build -mod=vendor .

lambda-build:
	rm -f test-api main.zip
	GOOS=linux go build
	zip main.zip test-api