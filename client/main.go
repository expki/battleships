package main

import (
	xencoding "battleships/src/encoding"
	"battleships/src/userinput"
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

	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Run game tick
	for i := 0; true; i++ {
		// Get bytes
		value := []byte(fmt.Sprintf("Hello, World! %d", i))
		data := xencoding.Encode(value)
		js.CopyBytesToJS(sharedArray, data)
		time.Sleep(time.Second)
	}
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
