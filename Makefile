.PHONY: build test clean install run

BINARY_NAME=finpup
INSTALL_PATH=/usr/local/bin

build:
	go build -o ${BINARY_NAME} cmd/finpup/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out

install: build
	sudo mv ${BINARY_NAME} ${INSTALL_PATH}/

run: build
	./${BINARY_NAME}

deps:
	go mod download
	go mod tidy
