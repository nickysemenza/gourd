FROM golang:1.20 AS builder-go

WORKDIR /work
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make bin/gourd

FROM rust:1.69 as builder-wasm
WORKDIR /work/rust
RUN cargo install wasm-pack
COPY rust/ .
WORKDIR /work
COPY Makefile .
RUN make generate-wasm

FROM node as build-ui
WORKDIR /work/ui
COPY ui/package.json ui/yarn.lock ./
RUN yarn
COPY --from=builder-wasm /work/rust/wasm/pkg /work/ui/src/wasm
COPY ui ./
RUN yarn build

FROM ghcr.io/nickysemenza/docker-magick:main
WORKDIR /work
COPY --from=builder-go /work/bin ./bin
COPY --from=builder-go /work/internal/db/migrations ./internal/db/migrations
COPY --from=build-ui /work/ui/build ./ui/build
COPY --from=builder-wasm /work/rust/wasm/pkg /wasm-pkg
ENTRYPOINT ["./bin/gourd","server"]