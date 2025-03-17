.PHONY: build-prod build-local start-azure start build-auth start-auth

build-prod:
	@cd azure_functions_go && GOOS=linux GOARCH=amd64 go build -o handler handler.go

build-local:
	@cd azure_functions_go && go build -o handler handler.go

start:
	@cd auth_service && go run cmd/main.go
