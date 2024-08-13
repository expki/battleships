import './wasm_exec';
declare class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
}

self.onmessage = async (event: MessageEvent<any>) => {
    const go = new Go();
    const { wasmModule } = event.data;
    const result = await WebAssembly.instantiate(wasmModule, go.importObject);
    go.run(result.instance);
}
