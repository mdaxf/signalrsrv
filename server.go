package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	//	"strings"

	//	"net"
	"net/http"
	"os"
	"time"

	//	"github.com/rs/cors"

	kitlog "github.com/go-kit/log"

	"github.com/mdaxf/signalrsrv/middleware"
	//	"github.com/mdaxf/signalrsrv/public"
	"github.com/mdaxf/signalrsrv/signalr"
	//	signalrserver "github.com/mdaxf/signalrsrv/signalr-server"
)

type Config struct {
	Address string `json:"address"`
	Clients string `json:"clients"`
}

var IACMessageBusName = "/iacmessagebus"

func runHTTPServer(address string, hub signalr.HubInterface, clients string) {
	server, _ := signalr.NewServer(context.TODO(), signalr.SimpleHubFactory(hub),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false),
		signalr.KeepAliveInterval(10*time.Second), signalr.AllowOriginPatterns([]string{clients}),
		signalr.InsecureSkipVerify(true))

	signalr.AllowedClients = clients

	router := http.NewServeMux()

	server.MapHTTP(signalr.WithHTTPServeMux(router), IACMessageBusName)

	//	fmt.Printf("Serving public content from the embedded filesystem\n")
	//	router.Handle("/", http.FileServer(http.FS(public.FS)))
	fmt.Printf("Listening for websocket connections on http://%s\n", address)
	if err := http.ListenAndServe(address, middleware.LogRequests(router)); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func runHTTPClient(address string, receiver interface{}) error {
	c, err := signalr.NewClient(context.Background(), nil,
		signalr.WithReceiver(receiver),
		signalr.WithConnector(func() (signalr.Connection, error) {
			creationCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			return signalr.NewHTTPConnection(creationCtx, address)
		}),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stdout), false))
	if err != nil {
		return err
	}
	c.Start()
	fmt.Println("Client started")
	return nil
}

type receiver struct {
	signalr.Receiver
}

func (r *receiver) Receive(msg string) {
	fmt.Println(msg)
	// The silly client urges the server to end his connection after 10 seconds
	r.Server().Send("abort")
}

func main() {
	appconfig := "signalRconfig.json"
	data, err := ioutil.ReadFile(appconfig)
	if err != nil {
		fmt.Errorf("failed to read configuration file: %v", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Errorf("failed to parse configuration file: %v", err)
	}

	address := config.Address

	clients := config.Clients

	//	url := "http://" + address + IACMessageBusName
	hub := &IACMessageBus{}

	go runHTTPServer(address, hub, clients)
	<-time.After(time.Millisecond * 2)
	/*	go func() {
		fmt.Println(runHTTPClient(url, &receiver{}))
	}() */
	ch := make(chan struct{})
	<-ch
}
