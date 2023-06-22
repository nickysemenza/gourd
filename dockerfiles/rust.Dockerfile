FROM lukemathwalker/cargo-chef:latest-rust-bullseye AS chef
WORKDIR /app

FROM chef AS planner
COPY rust/ .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder-rs 
COPY --from=planner /app/recipe.json recipe.json
# Build dependencies - this is the caching Docker layer!
RUN cargo chef cook --release --recipe-path recipe.json
# RUN cargo install wasm-pack
# Build application
COPY rust/ .
RUN cargo build --release --bin gourd

# We do not need the Rust toolchain to run the binary!
FROM debian:bullseye-slim AS runtime
WORKDIR /app
COPY --from=builder-rs /app/target/release/gourd /usr/local/bin

ENV RUST_BACKTRACE 1
ENV RUST_LOG debug
ENTRYPOINT ["/usr/local/bin/gourd","server"]


