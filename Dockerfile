FROM golang:1.13 AS builder


# Copy the code from the host and compile it
WORKDIR /tmp

COPY ./ ./
RUN make build

FROM scratch
COPY --from=builder /app ./
EXPOSE 8080 
ENTRYPOINT ["./app"]