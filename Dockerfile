FROM golang:1.17.6-alpine3.15 AS build

COPY src /src
WORKDIR /src

RUN GOOS=js GOARCH=wasm go build -o kafkagosaur.wasm

FROM scratch AS export
COPY --from=build /src/kafkagosaur.wasm /
