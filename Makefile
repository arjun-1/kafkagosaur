all: build

build:
	cd src; GOOS=js GOARCH=wasm go build -o ../bin/wasm.wasm

run:
	deno run --allow-read --allow-net lib/index.ts

test:
	docker-compose up -d

clean-test:
	docker-compose down

clean-build:
	rm bin/wasm.wasm

lint:
	deno fmt
	gofmt -w .