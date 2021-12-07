FROM alpine:3.13

COPY ./main /
COPY ./public /public
COPY ./config.toml /config.toml

CMD ["/main"]
