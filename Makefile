BINARY_NAME=ai-telegram-bot
DOCKER_IMAGE=ai-telegram-bot-image

GOPATH=$(shell go env GOPATH)

COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL := build

build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} ./internal/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@go clean

docker-build:
	@echo "Building Docker image..."
	docker build -t ${DOCKER_IMAGE} .

docker-start:
	docker run -d ${DOCKER_IMAGE}

docker-stop:
	docker stop $$(docker ps -q --filter ancestor=${DOCKER_IMAGE})

run: build
	./${BINARY_NAME}
