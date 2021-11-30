.PHONY: test

all: build

build:
	cd src; GOOS=js GOARCH=wasm go build -o ../bin/kafkagosaur.wasm

run:
	deno run --allow-read --allow-net lib/index.ts

test:
	docker-compose up -d

setup-test-data:


clean-test:
	docker-compose down

clean-build:
	rm bin/kafkagosaur.wasm

lint:
	deno fmt
	gofmt -w .