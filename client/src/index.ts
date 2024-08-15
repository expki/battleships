import { decode } from './encoding';
import * as enums from './enums';
import type * as types from './types';

if (!crossOriginIsolated) {
    console.error("sharredArrayBuffer is not available");
}

const canvas: HTMLCanvasElement = <HTMLCanvasElement>document.getElementById('main');
const ctx: CanvasRenderingContext2D = canvas.getContext("2d");

// Create a SharedArrayBuffer
const sharedBuffer = new SharedArrayBuffer(1024); // 1KB buffer

async function main(ev: Event) {
    // Check if canvas is supported
    if (!ctx) {
        console.error("Canvas not supported");
        return;
    }
    // Set canvas size to window size
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    // Set canvas style
    ctx.imageSmoothingEnabled= false;
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // Start game logic
    const logic = new Worker(new URL("./logic.ts", import.meta.url), { type: "module" });
    logic.onmessage = (ev) => {
        console.log("message:", ev.data);
    }
    const bytes = await fetch('logic.wasm').then(response => response.arrayBuffer());
    const payload: types.Payload<types.PayloadWasm> = {
        kind: enums.PayloadKind.wasm,
        payload: {
            wasm: bytes,
            pipe: sharedBuffer,
        },
    };
    logic.postMessage(payload)
    // Register input events
    const keyState: Record<string, boolean> = {};
    let mouseLeftState: boolean | undefined;
    let xState: number | undefined;
    let yState: number | undefined;
    addEventListener("keydown", (ev) => {
        if (keyState[ev.key] === true) return;
        keyState[ev.key] = true;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                keydown: ev.key,
            },
        };
        logic.postMessage(payload);
    });
    addEventListener("keyup", (ev) => {
        if (keyState[ev.key] === false) return;
        keyState[ev.key] = false;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                keyup: ev.key,
            },
        };
        logic.postMessage(payload);
    });
    addEventListener("mousedown", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== true) {
            enabled = true;
            data.mouseleft = true;
        }
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mouseup", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.button === 0 && mouseLeftState !== false) {
            enabled = true;
            data.mouseleft = false;
        }
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: data,
        };
        logic.postMessage(payload);
    });
    addEventListener("mousemove", (ev) => {
        let enabled = false;
        const data: types.PayloadInput = {};
        if (ev.clientX !== xState) {
            enabled = true;
            data.mousex = ev.clientX;
        }
        if (ev.clientY !== yState) {
            enabled = true;
            data.mousey = ev.clientY;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                mousex: ev.clientX,
                mousey: ev.clientY,
            },
        };
        logic.postMessage(payload);
    });
    const resizeScreen = (width: number, height: number): void => {
        let enabled = false;
        if (canvas.width !== width) {
            enabled = true;
            canvas.width = width; 
        }
        if (canvas.height !== height) {
            enabled = true;
            canvas.height = height;
        }
        if (!enabled) return;
        const payload: types.Payload<types.PayloadInput> = {
            kind: enums.PayloadKind.input,
            payload: {
                width: width,
                height: height,
            },
        };
        logic.postMessage(payload);
    }
    let resizeTimeout: NodeJS.Timeout | undefined;
    addEventListener("resize", () => {
        clearTimeout(resizeTimeout);
        // resize with backoff
        resizeTimeout = setTimeout(() => {
            resizeScreen(window.innerWidth, window.innerHeight);
        }, 500);
    });
    // Set initial values
    const intial: types.Payload<types.PayloadInput> = {
        kind: enums.PayloadKind.input,
        payload: {
            width: window.innerWidth,
            height: window.innerHeight,
        },
    };
    logic.postMessage(intial);
    // Start render loop
    renderLoop();
}

async function renderLoop(): Promise<void> {
    // Clear canvas
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // Draw
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, 100, 100);
    // Request next frame
    requestAnimationFrame(renderLoop);
}

addEventListener("load", main);

setInterval(() => {
    const sharedBytes = new Uint8Array(sharedBuffer);
    const bytes = new Uint8Array(sharedBytes.length);
    bytes.set(sharedBytes);
    const data = decode(bytes)
    console.log(data);
}, 1000);
