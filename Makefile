include .env
export

build-output:
	@mkdir -p bin/

build-server: cmd/server/main.go build-output
	@cd cmd/server/ && go build -o ../../bin/server

build-client: cmd/client/main.go build-output
	@cd cmd/client/ && go build -o ../../bin/client

run-build-server: build-server
	@./bin/server

run-build-client: build-client
	@./bin/client

run-server: cmd/server/main.go
	@go run cmd/server/*

run-client: cmd/client/main.go
	@go run cmd/client/*

build-win: cmd/server/main.go cmd/server/main.go
	@cd cmd/server && GOOS=windows GOARCH=amd64 go build -o ../../bin/server.exe
	@cd cmd/client && GOOS=windows GOARCH=amd64 go build -o ../../bin/client.exe

tst: test/files/test.go
	@go run test/files/test.go
