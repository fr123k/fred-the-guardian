.PHONY: build

go-init:
	go mod init github.com/fr123k/fred-the-guardian
	go mod vendor

build:
	go build -o build/main cmd/main.go
	go test -v -timeout 60s --cover -coverprofile=./build/cover.tmp ./...

coverage: build
	go tool cover -html=build/cover.tmp

run: build
	./build/main

clean:
	rm -rfv ./build
