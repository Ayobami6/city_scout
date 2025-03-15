.PHONY: build-prod build-local start

build-prod:
	@GOOS=linux GOARCH=amd64 go build handler.go

build-local:
	@go build handler.go

start:
	@go run handler.go