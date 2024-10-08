import { listenForUserInput } from './userinput/listener';
import { game } from './game/game';
import * as enums from './types/enums';
import type * as types from './types/types';

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
    const logic = new Worker(new URL("./logic/logic.ts", import.meta.url), { type: "module" });
    const bytes = await fetch('logic.wasm').then(response => response.arrayBuffer());
    const payload: types.Payload<types.PayloadWasm> = {
        kind: enums.PayloadKind.wasm,
        payload: {
            wasm: bytes,
            pipe: sharedBuffer,
        },
    };
    logic.postMessage(payload)
    // Listen for user input
    listenForUserInput(canvas, logic);
    // Set initial window size values
    const intial: types.Payload<types.PayloadInput> = {
        kind: enums.PayloadKind.input,
        payload: {
            width: window.innerWidth,
            height: window.innerHeight,
        },
    };
    logic.postMessage(intial);
    // Start render loop
    game(ctx, canvas, sharedBuffer);
}

addEventListener("load", main);
