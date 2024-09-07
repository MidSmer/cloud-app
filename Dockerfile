# go build
FROM golang:1.22-alpine3.20 as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# web pack
FROM node:20 as webpack

WORKDIR /build

COPY ./web/react-app ./web/react-app

RUN cd ./web/react-app \
    && npm install \
    && npm run build

# deploy
FROM alpine:3.20 as deploy

WORKDIR /root

COPY ./config.toml .

COPY --from=builder /build/app .
COPY --from=webpack /build/public ./public

CMD ["./app"]
