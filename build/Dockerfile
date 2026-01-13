FROM alpine

WORKDIR /app

RUN addgroup -g 1337 proxy && adduser -D -u 1337 -G proxy proxy

USER 1337

COPY _output/sidecar /app/sidecar
