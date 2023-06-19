FROM rust:1.70 as builder-wasm
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
ENV NODE_OPTIONS=--max_old_space_size=4096
COPY ui ./
RUN yarn build

# can be tiny image, just intermediary
FROM alpine:3 AS tmp
COPY --from=build-ui /work/ui/dist /work/ui/dist 
COPY --from=builder-wasm /work/rust/wasm/pkg /work/ui/src/wasm