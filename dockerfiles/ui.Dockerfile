FROM lukemathwalker/cargo-chef:latest-rust-bullseye AS chef
WORKDIR /work/rust

FROM chef AS planner
COPY rust/ .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder-rs 
COPY --from=planner /work/rust/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json
RUN cargo install wasm-pack
COPY rust/ .
RUN wasm-pack build wasm


FROM node as build-ui
WORKDIR /work/ui
COPY ui/package.json ui/yarn.lock ./
RUN yarn
COPY --from=builder-rs /work/rust/wasm/pkg /work/ui/src/wasm
ENV NODE_OPTIONS=--max_old_space_size=4096
COPY ui ./
RUN yarn build

# can be tiny image, just intermediary
FROM alpine:3 AS tmp
COPY --from=build-ui /work/ui/dist /work/ui/dist 
COPY --from=builder-rs /work/rust/wasm/pkg /work/ui/src/wasm