FROM rust:1.68 as builder
WORKDIR /work/rust
RUN apt-get update
COPY rust/ .
WORKDIR /work/rust/s
RUN pwd && ls
RUN cargo build
RUN ls
ENV RUST_BACKTRACE 1

ENTRYPOINT ["/work/rust/target/debug/gourd","server"]
