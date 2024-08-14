export type Payload<T> = {
    kind: number, 
    payload: T,
};

export type PayloadWasm = ArrayBuffer;

export type PayloadInput = {
    width?: number,
    height?: number,
    keyup?: string,
    keydown?: string,
    mouseleft?: boolean,
    mouseEpoch?: number,
    mousex?: number,
    mousey?: number,
};
