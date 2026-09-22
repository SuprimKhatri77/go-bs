.PHONY: fmt vet test race check generate

fmt:
	gofmt -l .

vet:
	go vet ./...

test:
	go test ./...

race:
	go test -race ./...

check: fmt vet test race

generate:
	go run ./tools/calendar-generator -out data.go
