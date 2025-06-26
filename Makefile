build:
	@go build -o bin/k8s-demo cmd/server/main.go

run: build
	@./bin/k8s-demo

deps:
	@go mod download
	@go mod verify

tidy:
	@go mod tidy