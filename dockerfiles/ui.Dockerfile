ARG IMAGE_TAG
FROM ghcr.io/nickysemenza/gourd-wasm-builder:${IMAGE_TAG} AS wasm

FROM node as build-ui
WORKDIR /work/ui
COPY ui/package.json ui/yarn.lock ./
RUN yarn
COPY --from=wasm /wasm /work/ui/src/wasm
ENV NODE_OPTIONS=--max_old_space_size=4096
COPY ui ./
RUN yarn build

# can be tiny image, just intermediary
FROM scratch AS tmp
COPY --from=build-ui /work/ui/dist /work/ui/dist 