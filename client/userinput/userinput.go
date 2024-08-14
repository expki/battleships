package userinput

import (
	"fmt"
	"sync"
	"syscall/js"
)

type UserInput struct {
	lock         sync.RWMutex
	windowWidth  int
	windowHeight int
}

func New() *UserInput {
	u := new(UserInput)
	window := js.Global()
	window.Set("handleInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		u.lock.Lock()
		defer u.lock.Unlock()
		fmt.Println(len(args))
		return nil
	}))
	return u
}

func (u *UserInput) Close() {
	u.lock.Lock()
	defer u.lock.Unlock()
}
