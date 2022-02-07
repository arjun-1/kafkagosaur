.PHONY: test

all: build-wasm

build-wasm:
	DOCKER_BUILDKIT=1 docker build -o bin .

test:
	deno test --allow-read --allow-net --unstable --coverage=coverage

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	deno fmt
	deno lint --ignore=lib/wasm_exec.js
	gofmt -w .

lint-ci:
	test -z $$(gofmt -l .)
	deno fmt --check --quiet
	deno lint --ignore=lib/wasm_exec.js --quiet
