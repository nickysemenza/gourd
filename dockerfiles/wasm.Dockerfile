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

FROM alpine AS tmp
COPY --from=builder-rs /work/rust/wasm/pkg /wasm