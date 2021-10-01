FROM rust:1.55 as builder
ENV SQLX_OFFLINE true
WORKDIR /work/rust
RUN apt-get update
RUN apt-get install -y pip libpython3.9-dev
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY rust/ .
WORKDIR /work/rust/s
RUN pwd && ls
RUN cargo build
RUN ls
ENV RUST_BACKTRACE 1
ENTRYPOINT ["/work/rust/target/debug/gourd"]