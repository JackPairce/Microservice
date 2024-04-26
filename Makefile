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

