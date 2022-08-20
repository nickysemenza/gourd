FROM golang:1.19 AS builder-go

WORKDIR /work
COPY . .
RUN go mod vendor
RUN make bin/gourd

FROM rust:1.63 as builder-wasm
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

FROM ghcr.io/nickysemenza/gourd-base:dev
RUN which pdflatex
RUN which magick
WORKDIR /work
COPY --from=builder-go /work/bin ./bin
COPY --from=builder-go /work/db/migrations ./db/migrations
COPY --from=build-ui /work/ui/build ./ui/build
ENTRYPOINT ["./bin/gourd","server"]