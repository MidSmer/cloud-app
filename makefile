all: main install build-web

main:
	CGO_ENABLED=0 go build -o main

install:
	cd web; \
	yarn install

build-web:
	cd web; \
	yarn build

clean:
	rm main
	rm -rf deploy
	rm -rf public

.PHONY: all main install build-web clean
