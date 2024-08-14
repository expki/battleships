import './wasm_exec';
import * as enums from './enums';
import type * as types from './types';

declare class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): void;
}

declare function handleInput(payload: types.PayloadInput): void;

self.onmessage = async (event: MessageEvent<types.Payload<any>>) => {
    switch (event.data.kind) {
        case enums.PayloadKind.wasm: {
            const payload: types.PayloadWasm = event.data.payload;
            const go = new Go();
            const result = await WebAssembly.instantiate(payload, go.importObject);
            go.run(result.instance);
            return;
        }            
        case enums.PayloadKind.input: {
            const payload: types.PayloadInput = event.data.payload;
            handleInput(payload);
            return;
        }
        default:
            console.error("Unknown payload kind:", event.data.kind);
            return;
    }
}
