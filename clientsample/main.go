package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mdaxf/signalrsrv/client"
)

func main() {
	// Create a client with the given address and hub name.
	client := client.NewWebsocketClient()

	// Set a function to be called when a server method is called.
	client.OnClientMethod = func(hub, method string, arguments []json.RawMessage) {
		fmt.Println("Message Received: ")
		fmt.Println("HUB: ", hub)
		fmt.Println("METHOD: ", method)
		fmt.Println("ARGUMENTS: ", arguments)
	}
	client.OnMessageError = func(err error) {
		fmt.Println("ERROR OCCURRED: ", err)
	}

	err := client.Connect("http", "127.0.0.1:8222", []string{"iacmessagebus"}) //and so forth

	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}
	defer client.Close()
	count := 0
	go func() {

		for count < 10 {
			client.CallHub("iacmessagebus", "send", "Test", "this is a message from the GO client")

			time.Sleep(5 * time.Second)
			count++
		}
	}()
}
