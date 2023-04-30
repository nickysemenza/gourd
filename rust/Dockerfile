FROM rust:1.69 as builder
WORKDIR /work/rust
RUN apt-get update
COPY rust/ .
WORKDIR /work/rust/s
RUN pwd && ls
RUN cargo build --release
RUN ls
ENV RUST_BACKTRACE 1
ENV RUST_LOG debug

ENTRYPOINT ["/work/rust/target/release/gourd","server"]
