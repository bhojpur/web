all: build

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o ../bin/webdbg

clean:
	@go clean
	@-rm web/app.wasm

run: build
	../bin/webdbg
