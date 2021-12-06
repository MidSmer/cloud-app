all: main setting install build-web

main:
	CGO_ENABLED=0 go build -o deploy/main

setting:
	cp config.toml deploy/

install:
	cd web; \
	yarn install

build-web:
	cd web; \
	yarn build

clean:
	rm -rf deploy

.PHONY: all main install build-web clean
