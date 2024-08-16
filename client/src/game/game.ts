import { Decode } from '../encoding/decoder';
import type GameState from '../types/state';

let ctx: CanvasRenderingContext2D;
let canvas: HTMLCanvasElement;
let sharedBuffer: SharedArrayBuffer;

export function game(inCtx: CanvasRenderingContext2D, inCanvas: HTMLCanvasElement, inSharedBuffer: SharedArrayBuffer) {
    ctx = inCtx;
    canvas = inCanvas;
    sharedBuffer = inSharedBuffer;
    renderLoop();
}

async function renderLoop(): Promise<void> {
    // Load game state
    const sharedBytes = new Uint8Array(sharedBuffer);
    const bytes = new Uint8Array(sharedBytes.length);
    bytes.set(sharedBytes);
    const state = Decode<GameState>(bytes)
    // Clear canvas
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // Draw
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, 100, 100);
    // Request next frame
    requestAnimationFrame(renderLoop);
}
