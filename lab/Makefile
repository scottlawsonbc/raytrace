.PHONY: client
client: clean
	GOOS=js GOARCH=wasm go build -o ./dist/client.wasm ./client
	cp ./static/favicon.ico ./dist/favicon.ico
	cp ./static/index.html ./dist/index.html
	cp -r ./assets ./dist/assets
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./dist/wasm_exec.js
	go build .
	go run .


# Auto restart server when any .go file is changed.
# Install with `go install github.com/cespare/reflex@latest`.
.PHONY: develop
develop:
	reflex -g '*/*.go' -s make

clean:
	rm -r ./dist
	mkdir ./dist
