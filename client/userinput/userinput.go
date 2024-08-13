package userinput

import (
	"sync"
	"syscall/js"
	"time"
)

type UserInput struct {
	lock           sync.RWMutex
	resizeHandler  js.Func
	keyDownHandler js.Func
	keyUpHandler   js.Func
	leftMouseDown  js.Func
	leftMouseUp    js.Func
	rightMouseDown js.Func
	rightMouseUp   js.Func
	windowWidth    int
	windowHeight   int
}

func New() *UserInput {
	window := js.Global().Get("window")
	u := new(UserInput)
	u.resizeHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		u.setWindowWidth(window.Get("innerWidth").Int())
		u.setWindowHeight(window.Get("innerHeight").Int())
		go func() {
			// Firefox sometimes doesn't update the window size immediately
			time.Sleep(500 * time.Millisecond)
			u.setWindowWidth(window.Get("innerWidth").Int())
			u.setWindowHeight(window.Get("innerHeight").Int())
		}()
		println("UserInput: Window resized:", u.GetWindowWidth(), u.GetWindowHeight())
		return nil
	})
	u.keyDownHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		key := event.Get("key").String()
		println("UserInput: Key down:", key)
		return nil
	})
	u.keyUpHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		key := event.Get("key").String()
		println("UserInput: Key up:", key)
		return nil
	})
	u.leftMouseDown = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		println("UserInput: Left mouse down:", event.Get("x").Int(), event.Get("y").Int())
		return nil
	})
	u.leftMouseUp = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		println("UserInput: Left mouse up:", event.Get("x").Int(), event.Get("y").Int())
		return nil
	})
	u.rightMouseDown = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		//event.Call("preventDefault") activate once ready to implement right click menu
		println("UserInput: Right mouse down:", event.Get("x").Int(), event.Get("y").Int())
		return nil
	})
	u.rightMouseUp = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		event := args[0]
		println("UserInput: Right mouse up:", event.Get("x").Int(), event.Get("y").Int())
		return nil
	})
	js.Global().Get("window").Call("addEventListener", "resize", u.resizeHandler)
	js.Global().Get("document").Call("addEventListener", "keydown", u.keyDownHandler)
	js.Global().Get("document").Call("addEventListener", "keyup", u.keyUpHandler)
	js.Global().Get("document").Call("addEventListener", "mousedown", u.leftMouseDown)
	js.Global().Get("document").Call("addEventListener", "mouseup", u.leftMouseUp)
	js.Global().Get("document").Call("addEventListener", "contextmenu", u.rightMouseDown)
	js.Global().Get("document").Call("addEventListener", "contextmenu", u.rightMouseUp)
	return u
}

func (u *UserInput) Close() {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.resizeHandler.Release()
	u.keyDownHandler.Release()
	u.keyUpHandler.Release()
	u.leftMouseDown.Release()
	u.leftMouseUp.Release()
	u.rightMouseDown.Release()
	u.rightMouseUp.Release()
}
