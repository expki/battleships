// Copyright [2019] [Mark Farnan]

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// NOTICE:  Much of this demo is a re-write of the 'Moving red Laser' demo by Martin Olsansky https://medium.freecodecamp.org/webassembly-with-golang-is-fun-b243c0e34f02
// It has been re-written to make use of the go-canvas library,  and avoid context calls for drawing.

package main

import (
	"battleships/userinput"
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"syscall/js"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var port string

func websocketURL() (string, error) {
	uri, err := url.Parse(js.Global().Get("location").Get("href").String())
	if err != nil {
		return "", fmt.Errorf("url.Parse window href: %v", err)
	}
	if port != "" {
		wsPort, err := strconv.Atoi(port)
		if err != nil {
			return "", fmt.Errorf("strconv.Atoi alt port: %v", err)
		}
		uri.Host = fmt.Sprintf("%s:%d", uri.Hostname(), wsPort)
	}
	switch uri.Scheme {
	case "http":
		uri.Scheme = "ws"
	case "https":
		uri.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported scheme: %q", uri.Scheme)
	}
	return uri.String(), nil
}

type Websocket struct {
	*websocket.Conn
}

func main() {
	// Listen to user input
	input := userinput.New()
	defer input.Close()
	select {}

	// Define the websocket event handler
	uri, err := websocketURL()
	if err != nil {
		log.Fatalf("websocketURL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, uri, nil)
	if err != nil {
		log.Fatalf("websocket.Dial: %v", err)
	}
	defer c.CloseNow()

	err = wsjson.Write(ctx, c, "hi")
	if err != nil {
		log.Fatalf("wsjson.Write: %v", err)
	}

	c.Close(websocket.StatusNormalClosure, "")

	select {}
}

func getFrameTime() time.Duration {
	const count = 20
	var n uint8
	var start time.Time
	done := make(chan struct{})
	var fn js.Func
	fn = js.FuncOf(func(this js.Value, args []js.Value) (void any) {
		if n == 0 {
			start = time.Now()
		} else if n >= count {
			close(done)
			return
		}
		js.Global().Call("requestAnimationFrame", fn)
		n++
		return
	})
	js.Global().Call("requestAnimationFrame", fn)
	<-done
	fn.Release()
	return time.Since(start) / count
}
