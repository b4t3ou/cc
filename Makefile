deps:
	go mod verify

lint:
	golangci-lint run

test:
	go test -cover -failfast ./...

