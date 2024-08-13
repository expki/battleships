const canvas: HTMLCanvasElement = <HTMLCanvasElement>document.getElementById('main');
const ctx: CanvasRenderingContext2D = canvas.getContext("2d");

function main(ev: Event) {
    // Check if canvas is supported
    if (!ctx) {
        console.error("Canvas not supported");
        return;
    }
    // Set canvas size to window size
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    addEventListener("resize", () => {
        resizeScreen();
        // sometimes firefox leaves a white border, this prevents it
        setTimeout(() => {
            resizeScreen();
        }, 500);
    });
    // Set canvas style
    ctx.imageSmoothingEnabled= false;
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // Start game logic
    const logic = new Worker(new URL("./logic.ts", import.meta.url), { type: "module" });
    logic.onmessage = (ev) => {
        console.log(ev);
    }
    fetch('logic.wasm').then(response =>
        response.arrayBuffer()
    ).then(bytes => {
        logic.postMessage({ payload: bytes });
    });
    // Start render loop
    renderLoop();
}

function resizeScreen() {
    if (canvas.height !== window.innerHeight) canvas.height = window.innerHeight;
    if (canvas.width !== window.innerWidth) canvas.width = window.innerWidth; 
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
