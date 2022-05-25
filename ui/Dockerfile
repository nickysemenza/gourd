FROM rust:1.61 as builder
WORKDIR /work/rust
RUN curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
COPY rust/ .
WORKDIR /work/rust
RUN pwd && ls
RUN wasm-pack build w

# Stage 1 - the build process
FROM node:buster-slim as build-deps
WORKDIR /work/ui
COPY ui/package.json ui/yarn.lock ./
RUN yarn
COPY --from=builder /work/rust/w/pkg /work/ui/src/wasm
COPY ui ./
RUN yarn build

# Stage 2 - the production environment
FROM nginx:1.22-alpine
COPY --from=build-deps /work/ui/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]