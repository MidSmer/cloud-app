all: main install build-web

main:
	CGO_ENABLED=0 go build -o main

install:
	cd web/react-app; \
	npm install

build-web:
	cd web/react-app; \
	npm run build

clean:
	rm main
	rm -rf deploy
	rm -rf public

.PHONY: all main install build-web clean
