include .env
export

build-output:
	@mkdir -p bin/

build-server: cmd/server/main.go build-output
	@cd cmd/server/ && go build -o ../../bin/server

build-client: cmd/client/main.go build-output
	@cd cmd/client/ && go build -o ../../bin/client

serve: build-server
	@./bin/server

cli: build-client
	@./bin/client

build-win: cmd/server/main.go cmd/server/main.go
	@cd cmd/server && GOOS=windows GOARCH=amd64 go build -o ../../bin/server.exe
	@cd cmd/client && GOOS=windows GOARCH=amd64 go build -o ../../bin/client.exe

tst: test/files/test.go
	@go run test/files/test.go
