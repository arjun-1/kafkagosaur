.PHONY: test

all: build

build:
	cd src; GOOS=js GOARCH=wasm go build -o ../bin/kafkagosaur.wasm

test:
	deno test --allow-read --allow-net

docker:
	docker-compose up

docker-clean:
	docker-compose down

lint:
	deno fmt
	gofmt -w .