add.wasm: add.go
	GOOS=wasip1 GOARCH=wasm tinygo build -o $@ $<

run: add.wasm main.go
	go run main.go
