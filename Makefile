fmt:
	go fmt ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run

check: fmt vet lint test