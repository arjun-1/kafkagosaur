all: build

build:
	GOOS=js GOARCH=wasm go build -o wasm.wasm

run:
	deno run --allow-read --allow-net index.ts

test:
	docker-compose up -d

clean-test:
	docker-compose down

clean-build:
	rm wasm.wasm