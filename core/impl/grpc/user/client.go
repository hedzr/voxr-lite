/*
 * Copyright © 2019 Hedzr Yeh.
 */

package user

// sample client for testing and verifying

import (
	"context"
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
	chLogin         chan *v10.LoginReq
}

type Client struct {
	conn          *grpc.ClientConn
	serviceClient v10.UserActionClient
	friendClient  v10.FriendActionClient
	callback      func(res *v10.UserInfoToken)
}

var hub = &Hub{
	clients:         make(map[*Client]bool),
	exitCh:          make(chan bool),
	exited:          true,
	newClientAdding: make(chan *Client),
	destroyClient:   make(chan *Client),
	chLogin:         make(chan *v10.LoginReq),
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

		case req := <-hub.chLogin:
			// logrus.Println(name)
			// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*len(hub.clients))*time.Second)
			// defer cancel()
			ctx := context.Background()
			for k, ok := range hub.clients {
				if ok && k.conn != nil {
					r, err := k.serviceClient.Login(ctx, req)
					if err != nil {
						logrus.Warnf("could not greet: %v", err)
					} else {
						logrus.Printf("Login: %s", r)
						if k.callback != nil {
							k.callback(r)
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

func NewClient(callback func(res *v10.UserInfoToken)) (client *Client) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(tool.LoadHostDefinition("server.deps.user"), grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()

	c := v10.NewUserActionClient(conn)
	f := v10.NewFriendActionClient(conn)
	client = &Client{conn: conn, serviceClient: c, friendClient: f, callback: callback}

	// notify hub routine the new client connection object is incoming
	hub.newClientAdding <- client
	return
}

func (c *Client) Login(req *v10.LoginReq, callback func(res *v10.UserInfoToken)) {
	// Contact the server and print out its response.
	// hub.chLogin <- req
}

func ClientSendUnique(req *v10.LoginReq, callback func(res *v10.UserInfoToken)) {
	if len(hub.clients) == 0 {
		noClose(req, callback)
	} else {
		go func() {
			sent := false
			for client, exists := range hub.clients {
				if exists {
					client.Login(req, callback)
					sent = true
				}
			}
			if sent == false {
				noClose(req, callback)
			}
		}()
	}
}

func ClientSend(req *v10.LoginReq, callback func(res *v10.UserInfoToken)) {
	autoClose(req, callback)
}

func autoClose(req *v10.LoginReq, callback func(res *v10.UserInfoToken)) {
	go func() {
		client := NewClient(callback)
		defer client.Close()
		time.Sleep(1 * time.Second)
		client.Login(req, callback)
	}()
}

func noClose(req *v10.LoginReq, callback func(res *v10.UserInfoToken)) {
	go func() {
		client := NewClient(callback)
		time.Sleep(1 * time.Second)
		client.Login(req, callback)
	}()
}
