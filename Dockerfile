FROM golang:1.17.5 AS builder


# Copy the code from the host and compile it
WORKDIR /work
COPY go.mod . 
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN go mod vendor
RUN make bin/gourd

FROM ghcr.io/nickysemenza/gourd-base:dev
RUN which pdflatex
RUN which magick
WORKDIR /work
COPY --from=builder /work/bin ./bin
COPY --from=builder /work/db/migrations ./db/migrations
ENTRYPOINT ["./bin/gourd","server"]