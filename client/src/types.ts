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
    mousex?: number,
    mousey?: number,
};
