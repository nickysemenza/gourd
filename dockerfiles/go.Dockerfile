FROM golang:1.20 AS builder-go

WORKDIR /work
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make bin/gourd

# can be tiny image, just intermediary
FROM alpine:3 AS tmp
WORKDIR /work
COPY --from=builder-go /work /work
COPY --from=builder-go /work/internal/db/migrations /work/internal/db/migrations