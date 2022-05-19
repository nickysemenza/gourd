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

# TODO: `configuration file "/work/s-configuration/base" not found`
# FROM debian:stable

# RUN apt-get update
# RUN apt-get install -y libpython3.9-dev
# WORKDIR /work
# COPY --from=builder /work/rust/target/debug ./bin
# ENTRYPOINT ["./bin/gourd","server"]
