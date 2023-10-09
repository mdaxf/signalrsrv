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
	client.OnErrorMethod = func (err error) {
	   fmt.Println("ERROR OCCURRED: ", err)
	}

	client.Connect("http", "127.0.0.1:8222", []string{"iacmessagebus"}) //and so forth

	defer client.Close()

	go func ()  {
		
		client.CallHub("iacmessagebus", "send", "Test", "this is a message from the GO client")
		
		time.Sleep(5 * time.Second)
	}()
}