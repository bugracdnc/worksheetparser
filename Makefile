.PHONY: build test web serve clean

build:
	go build -o worksheetparser ./cmd/worksheetparser

test:
	go test ./...

# Browser editor: web/ becomes a static site (GitHub Pages serves it as-is).
web:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/main.wasm ./cmd/wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
	cp examples/example.wsp.txt web/

serve: web
	python3 -m http.server -d web 8080

clean:
	rm -f worksheetparser web/main.wasm web/wasm_exec.js web/example.wsp.txt
