
serve: cmd/server/main.go
	@go run cmd/server/main.go

cli: cmd/client/main.go
	@go run cmd/client/*.go

