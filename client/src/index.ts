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
    // Start game loop
    gameLoop();
}

function resizeScreen() {
    if (canvas.height !== window.innerHeight) canvas.height = window.innerHeight;
    if (canvas.width !== window.innerWidth) canvas.width = window.innerWidth; 
}

async function gameLoop(): Promise<void> {
    // Clear canvas
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    // Draw
    ctx.fillStyle = "white";
    ctx.fillRect(0, 0, 100, 100);
    // Request next frame
    requestAnimationFrame(gameLoop);
}

addEventListener("load", main);
