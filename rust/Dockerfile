FROM rust:1.61 as builder
WORKDIR /work/rust
RUN apt-get update
RUN apt-get install -y pip libpython3.9-dev
COPY rust/ .
RUN pip install -r requirements.txt
WORKDIR /work/rust/s
RUN pwd && ls
RUN cargo build
RUN ls
ENV RUST_BACKTRACE 1

ENTRYPOINT ["/work/rust/target/debug/gourd","server"]
