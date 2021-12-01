.PHONY: build

VERSION=3.3
PORT?=8080
export NAME=fr123k/fred-the-guardian
export IMAGE="${NAME}:${VERSION}"
export LATEST="${NAME}:latest"

export DOCKER_COMMAND_LOCAL=docker run \
		-e PORT="${PORT}" \
		-p ${PORT}:${PORT} \

docker-build:
	docker build -t $(IMAGE) -f Dockerfile .
	
go-init:
	go mod init github.com/fr123k/fred-the-guardian
	go mod vendor

build:
	go build -o build/main srv/ping.go
	go test -v -timeout 60s --cover -coverprofile=./build/cover.tmp ./srv ./pkg/...

build-cli:
	go build -o build/pong cmd/pong.go
	go test -v -timeout 60s --cover -coverprofile=./build/cover.tmp ./cmd

coverage:
	go test -v -timeout 60s --cover -coverprofile=./build/cover.tmp ./...
	go tool cover -html=build/cover.tmp

build-all: build build-cli

run: build
	./build/main

docker-run: docker-build
	docker stop fred || echo ignore error
	$(DOCKER_COMMAND_LOCAL) -d --rm --name fred  $(IMAGE)

clean:
	docker stop fred || echo ignore error
	rm -rfv ./build

release: docker-build ## Push docker image to docker hub
	docker tag ${IMAGE} ${LATEST}
	docker push ${IMAGE}
	docker push ${NAME}

test: docker-run
	curl -X POST -H 'X-SECRET-KEY:top secret' -v http://localhost:${PORT}/ping -d '{"request":"ping"}'

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this help
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
