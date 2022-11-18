# go build
FROM golang:1.19-alpine3.15 as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# web pack
FROM node:18 as webpack

WORKDIR /build

COPY ./web ./web

RUN cd ./web \
    && yarn install \
    && yarn build

# deploy
FROM alpine:3.15 as deploy

WORKDIR /root

COPY ./config.toml .

COPY --from=builder /build/app .
COPY --from=webpack /build/public ./public

CMD ["./app"]
