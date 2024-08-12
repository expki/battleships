GOOS=js GOARCH=wasm go build -trimpath -ldflags="-X 'main.port=8001'" -o battleships.wasm .
python3 -m http.server
