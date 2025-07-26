gen-mock:
	go generate ./...

lint:
	golangci-lint run

