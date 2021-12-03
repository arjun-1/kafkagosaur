.PHONY: test

all: build

build:
	cd src; GOOS=js GOARCH=wasm go build -o ../bin/kafkagosaur.wasm

run:
	deno run --allow-read --allow-net lib/index.ts

test:
	deno test --allow-read --allow-net

docker:
	docker-compose up

lint:
	deno fmt
	gofmt -w .