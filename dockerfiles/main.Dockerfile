
ARG IMAGE_TAG

FROM ghcr.io/nickysemenza/gourd-ui-builder:${IMAGE_TAG} AS build-ui
FROM ghcr.io/nickysemenza/gourd-go-builder:${IMAGE_TAG} AS builder-go

FROM ghcr.io/nickysemenza/docker-magick:main
WORKDIR /work
COPY --from=builder-go /work/bin ./bin
COPY --from=builder-go /work/internal/db/migrations ./internal/db/migrations
COPY --from=build-ui /work/ui/build ./ui/build
ENTRYPOINT ["./bin/gourd","server"]