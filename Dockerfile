# syntax=docker/dockerfile:1

FROM node:24.13.0-trixie AS fe

WORKDIR /app
COPY . /app

RUN corepack enable
RUN corepack prepare --activate

RUN apt-get update && apt-get install -y \
    make \
    && rm -rf /var/lib/apt/lists/*

RUN make install build_frontend

FROM golang:1.25.7-trixie AS be

WORKDIR /app
COPY . /app

COPY --from=fe /app/packages/api/dist /app/packages/api/dist

RUN apt-get update && apt-get install -y \
    make \
    && rm -rf /var/lib/apt/lists/*

RUN CGO_ENABLED=0 make build_backend

FROM alpine:3.23.3

WORKDIR /app

COPY --from=be /app/build/bin/strshelf /app

RUN apk add --no-cache tzdata

CMD [ "/app/strshelf" ]
