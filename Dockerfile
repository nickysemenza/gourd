FROM golang:1.14 AS builder


# Copy the code from the host and compile it
WORKDIR /work
COPY go.mod . 
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN make bin/food

FROM debian:buster
WORKDIR /work
COPY --from=builder /work/bin ./bin
COPY --from=builder /work/migrations ./migrations
RUN ls
EXPOSE 4242
ENTRYPOINT ["./bin/food"]