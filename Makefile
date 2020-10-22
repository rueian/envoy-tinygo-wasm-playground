wasm:
	tinygo build -o /src/main.wasm -scheduler=none -target=wasi -wasm-abi=generic /src/main.go