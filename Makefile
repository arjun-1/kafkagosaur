.PHONY: test

all: build

build:
	docker build -o bin .

test:
	deno test --allow-read --allow-net --coverage=coverage

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
