/*
 * Copyright © 2019 Hedzr Yeh.
 */

package apply

// sample client for testing and verifying

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/tool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

const (
	defaultName = "world"
)

type Hub struct {
	clients         map[*Client]bool
	exitCh          chan bool
	exited          bool
	newClientAdding chan *Client
	destroyClient   chan *Client
	querying        chan []byte
}

type Client struct {
	conn               *grpc.ClientConn
	applyServiceClient v10.ApplyServiceClient
	callback           func(string)
}

var hub = &Hub{
	clients:         make(map[*Client]bool),
	exitCh:          make(chan bool),
	exited:          true,
	newClientAdding: make(chan *Client),
	destroyClient:   make(chan *Client),
	querying:        make(chan []byte),
}

func StartClient() {
	go hub.start()
}

func StopClient() {
	if !hub.exited {
		hub.exitCh <- true
	}
}

func (hub *Hub) start() {
	logrus.Println("grpc hub started.")
	for {
		hub.exited = false
		select {
		case exit := <-hub.exitCh:
			if exit {
				logrus.Infof("grpc hub exiting.")
				for k, ok := range hub.clients {
					if ok && k.conn != nil {
						k.conn.Close()
					}
				}
				hub.exited = exit
				return
			}

		case c := <-hub.newClientAdding:
			hub.clients[c] = true                                               // flag the client conn object is true inside map `clients`
			logrus.Println("=== new grpc client session created.", hub.clients) // and log it for debugging

		case c := <-hub.destroyClient:
			delete(hub.clients, c)
			c.conn.Close()

		case name := <-hub.querying:
			// logrus.Println(name)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*len(hub.clients))*time.Second)
			// defer cancel()
			ctx := context.Background()
			for k, ok := range hub.clients {
				if ok && k.conn != nil {
					r, err := k.applyServiceClient.FetchApplyByUid(ctx, &v10.ApplyRequest{Uid: string(name)})
					if err != nil {
						logrus.Warnf("could not greet: %v", err)
					} else if !r.Ok {
						logrus.Warnf("could not greet - server return failed: code=%d, %v", r.ErrCode, r.Msg)
					} else {
						logrus.Printf("Greeting: %v", r.Data)
						if k.callback != nil {
							ret := &v10.ApplyResponse{}
							if err := ptypes.UnmarshalAny(r.Data[0], ret); err == nil {
								k.callback(ret.Result)
							} else {
								logrus.Warnf("cannot decode to pb.ApplyResponse: %v", r.Data)
							}
						}
					}
				}
			}

		}
	}
}

func (h *Hub) Close(client *Client) {
	if enabled, ok := hub.clients[client]; ok && enabled {
		hub.destroyClient <- client
	}
}

func (c *Client) Close() {
	hub.Close(c)
}

func NewClient(callback func(string)) (client *Client) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(tool.LoadHostDefinition("server.deps.apply"), grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()

	c := v10.NewApplyServiceClient(conn)
	client = &Client{conn: conn, applyServiceClient: c, callback: callback}

	// notify hub routine the new client connection object is incoming
	hub.newClientAdding <- client
	return
}

func (c *Client) QueryByName(name string) {
	// Contact the server and print out its response.
	hub.querying <- []byte(name)
}

func ClientSend(name string, callback func(string)) {
	autoClose(name, callback)
}

func autoClose(name string, callback func(string)) {
	go func() {
		client := NewClient(callback)
		defer client.Close()
		time.Sleep(1 * time.Second)
		client.QueryByName(name)
	}()
}

func noClose(name string, callback func(string)) {
	go func() {
		client := NewClient(callback)
		time.Sleep(1 * time.Second)
		client.QueryByName(name)
	}()
}

func ClientSendUnique(name string, callback func(string)) {
	if len(hub.clients) == 0 {
		noClose(name, callback)
	} else {
		go func() {
			sent := false
			for client, exists := range hub.clients {
				if exists {
					client.QueryByName(name)
					sent = true
				}
			}
			if sent == false {
				noClose(name, callback)
			}
		}()
	}
}
