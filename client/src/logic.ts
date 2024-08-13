import './wasm_exec';
declare class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
}

type payload<T> = {
    kind: PayloadKind, 
    payload: T,
};

enum PayloadKind {
    wasm = 0,
    resize = 1,
    keyUp = 2,
    keyDown = 3,
    mouseUp = 4,
    contextUp = 6,
}

type PayloadWasm = ArrayBuffer;
type PayloadResize = {
    width: number,
    height: number,
};
type PayloadKeyUp = string;
type PayloadKeyDown = string;
type PayloadMouseUp = boolean;
type PayloadContextUp = boolean;

self.onmessage = async (event: MessageEvent<payload<any>>) => {
    switch (event.data.kind) {
        case PayloadKind.wasm: {
            const payload: PayloadWasm = event.data.payload;
            const go = new Go();
            const result = await WebAssembly.instantiate(payload, go.importObject);
            go.run(result.instance);
            return;
        }            
        case PayloadKind.resize: {
            const payload: PayloadResize = event.data.payload;
            return;
        }
        default:
            console.error("Unknown payload kind:", event.data.kind);
            return;
    }
}
