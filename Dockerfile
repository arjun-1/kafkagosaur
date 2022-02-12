# syntax=docker/dockerfile:1.3

FROM golang:1.17.6-alpine3.15 AS build

COPY src/go.* /src/
WORKDIR /src

RUN go mod download
COPY src /src/

RUN --mount=type=cache,target=/root/.cache/go-build \
GOOS=js GOARCH=wasm go build -o kafkagosaur.wasm

FROM scratch AS export

COPY --from=build /src/kafkagosaur.wasm /
