FROM golang:1.17 AS builder


# Copy the code from the host and compile it
WORKDIR /work
COPY go.mod . 
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN make bin/gourd

# FROM debian:buster
FROM alpine:20210804
RUN apk add --no-cache ca-certificates
# https://stackoverflow.com/a/35613430/1374045
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /work
COPY --from=builder /work/bin ./bin
COPY --from=builder /work/db/migrations ./db/migrations
RUN ls
ENTRYPOINT ["./bin/gourd","server"]