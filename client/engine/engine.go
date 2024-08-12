package engine

import (
	"fmt"
	"syscall/js"
	"time"
)

type Engine struct {
	window   js.Value
	document js.Value
	canvas   js.Value
	ctx      js.Value
}

func NewEngine() *Engine {
	engn := &Engine{
		window: js.Global(),
	}

	// Canvas
	engn.document = engn.window.Get("document")
	engn.canvas = engn.document.Call("createElement", "canvas")
	engn.ctx = engn.canvas.Call("getContext", "2d")

	// Append the canvas to the document body
	width := engn.window.Get("innerWidth").Int()
	height := engn.window.Get("innerHeight").Int()
	engn.canvas.Set("width", width)
	engn.canvas.Set("height", height)
	engn.document.Get("body").Call("appendChild", engn.canvas)
	fmt.Println("Canvas size:", width, height)

	// Loop
	n := int(0)
	start := time.Now()
	var callback js.Func
	callback = js.FuncOf(func(this js.Value, args []js.Value) (void any) {
		screen := make([][][3]int, height)
		for i := range screen {
			screen[i] = make([][3]int, width)
			for j := range screen[i] {
				screen[i][j] = [3]int{n % 50, 0, 0}
			}
		}
		engn.render(screen)
		n++
		if n%60 == 0 {
			fmt.Println("FPS:", 60/time.Since(start).Seconds())
			start = time.Now()
		}
		engn.window.Call("requestAnimationFrame", callback)
		return
	})
	engn.window.Call("requestAnimationFrame", callback)

	return engn
}

func (e *Engine) render(grid [][][3]int) {
	height := len(grid)
	if height == 0 {
		return
	}
	width := len(grid[0])

	// Create a new Uint8ClampedArray for ImageData
	// Each pixel requires 4 values (RGBA), hence the width * height * 4
	pixelData := make([]uint8, width*height*4)
	index := 0
	for _, row := range grid {
		for _, cell := range row {
			// Populate the pixel data, cell contains [r, g, b]
			pixelData[index] = uint8(cell[0])   // R
			pixelData[index+1] = uint8(cell[1]) // G
			pixelData[index+2] = uint8(cell[2]) // B
			pixelData[index+3] = 255            // A (fully opaque)
			index += 4
		}
	}

	// Convert the Go slice to a JavaScript TypedArray
	jsPixelData := e.window.Get("Uint8ClampedArray").New(len(pixelData))
	js.CopyBytesToJS(jsPixelData, pixelData)

	// Create ImageData from the pixel data and dimensions
	imageData := e.window.Get("ImageData").New(jsPixelData, width, height)

	// Put the ImageData onto the canvas
	e.ctx.Call("putImageData", imageData, 0, 0) // Adjust (0, 0) to your desired position
}
