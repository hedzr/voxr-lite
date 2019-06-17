/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type (
	WsHub struct {
		mutex   *sync.RWMutex
		clients map[*WsClient]bool
		// clientsMap map[string]*WsClient //
		exitCh       chan bool
		exited       bool
		broadcast    chan []byte
		register     chan *WsClient // request to add a new client
		unregister   chan *WsClient // request to delete a client
		addClient    chan bool
		removeClient chan *WsClient
		pullMessages chan *PullMessages
		// mesgReading chan *TxtMsg
		// mqttReading chan *BinMsg
		// ppttReading chan *BinMsg
		// ppttPushing chan *newMsgIncoming
		seq uint32
	}

	PullMessages struct {
		uid    uint64
		did    string
		notify *v10.NotifyMessage
	}
)

var wshub = &WsHub{
	mutex:   new(sync.RWMutex),
	clients: make(map[*WsClient]bool),
	// clientsMap: make(map[string]*WsClient),
	exitCh:       make(chan bool),
	exited:       true,
	broadcast:    make(chan []byte),
	register:     make(chan *WsClient),
	unregister:   make(chan *WsClient),
	addClient:    make(chan bool),
	removeClient: make(chan *WsClient),
	pullMessages: make(chan *PullMessages),
	// mesgReading: make(chan *TxtMsg),
	// mqttReading: make(chan *BinMsg),
	// ppttReading: make(chan *BinMsg),
	// ppttPushing: make(chan *newMsgIncoming),
}

func WsStart() {
	go wshub.start()
}

func WsStop() {
	wshub.stop()
}

func WsClientAddNew() {
	wshub.addClient <- true
}

func (h *WsHub) stop() {
	if !h.exited {
		h.exited = true
		if h.exitCh != nil {
			h.exitCh <- true
			time.Sleep(1 * time.Second)
			close(h.exitCh)
			h.exitCh = nil
		}

		for k, v := range h.clients {
			if v {
				k.CloseFully()
			}
		}
		h.clients = make(map[*WsClient]bool)

		if h.broadcast != nil {
			close(h.broadcast)
			h.broadcast = nil
		}
		if h.register != nil {
			close(h.register)
			h.register = nil
		}
		if h.unregister != nil {
			close(h.unregister)
			h.unregister = nil
		}
		if h.addClient != nil {
			close(h.addClient)
			h.addClient = nil
		}
		if h.removeClient != nil {
			close(h.removeClient)
			h.removeClient = nil
		}
	}
}

func (h *WsHub) addNewClient() {
	go func() {

		wsClient := &WsClient{
			ReconnectDuration: 3 * time.Second,
		}
		if err := wsClient.Init(); err == nil {
			// the client will be registered to ws-hub once while it's connected to remote.
			wsClient.Open()
		}

	}()
}

func (h *WsHub) start() {
	h.preInitCommands()
	loadMaxMessageSize()

	getTicker := time.NewTicker(13 * time.Second)
	sendTicker := time.NewTicker(7 * time.Second)
	ticker := time.NewTicker(11 * time.Second)
	defer func() {
		ticker.Stop()
		sendTicker.Stop()
		getTicker.Stop()
		logrus.Debugf("chat hub.run() stopped.")
	}()

	for {
		h.exited = false
		select {
		case exit := <-h.exitCh:
			if exit {
				logrus.Infof("    [ws] chat hub exiting. (WebSocket message processing service)")
				h.exited = true
				return
			}

		case client := <-h.register:
			h.onRegister(client)

		case client := <-h.unregister:
			logrus.Debugf("    [ws] unregistering ws client from hub: %v", client)
			h.onDeregister(client)

		case _ = <-h.addClient:
			h.addNewClient()

		// case binMsg := <-h.ppttReading:
		// 	h.ppttDoProcessMsg(binMsg)
		//
		// case msg := <-h.ppttPushing:
		// 	h.ppttDoPushing(msg)
		//
		// case users := <-service.UsersNeedNotified:
		// 	logrus.Debugf("those users need be notified: %v", *users)
		// 	h.onNotifyUsersNewMsgIncoming(users)

		case message := <-h.broadcast:
			h.doBroadcast(message)

		case po := <-h.pullMessages:
			if err := h.pullNewMessages(po); err != nil {
				logrus.Errorf("    [ws] ERR (pullMessages): %v", err)
				return
			}

		case tm := <-sendTicker.C:
			if vxconf.GetBoolR("server.websocket.send.send-msg", true) {
				if err := h.testForSendMsg(tm); err != nil {
					logrus.Errorf("    [ws] ERR (sendTicker): %v", err)
					return
				}
			}

		case tm := <-getTicker.C:
			if vxconf.GetBoolR("server.websocket.send.get-off-msg", false) {
				if err := h.testForGetOffline(tm); err != nil {
					logrus.Errorf("    [ws] ERR (sendTicker): %v", err)
				}
			}

		case tm := <-ticker.C:
			if vxconf.GetBoolR("server.websocket.log.dump-clients", false) {
				logrus.Debugf("--- [ws] ws clients at %v: ", tm)
				h.mutex.RLock()
				for client := range h.clients {
					logrus.Debugf("      [ws] clients: %v", client)
				}
				h.mutex.RUnlock()
			} else {
				logrus.Debugf("--- [ws] ws clients at %v: %v clients", tm, len(h.clients))
			}
		}
	}
}

func (h *WsHub) doBroadcast(message []byte) {
	logrus.Debugf("    [ws] broadcasting message: %v", message)
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for client := range h.clients {
		client.sendTxtMsg(message)
		// select {
		// case client.textSend <- message:
		// default:
		// 	//close(client.textSend)
		// 	//delete(h.clients, client)
		// }
	}
}

func (h *WsHub) preInitCommands() {
	// // apps/clients 通过 websocket 发送 PB 消息，Hub 负责将这些消息通过 Hub.commands 分发到具体的处理逻辑
	// h.commands = make(map[v10.Op]func(from *WsClient, body []byte))
	// h.commands[v10.Op_SendMsg] = h.cxSendMsg
}

func (h *WsHub) onRegister(client *WsClient) {
	// h.mutex.RLock()
	// if ok, _ := h.clients[client]; !ok {
	// 	h.mutex.RUnlock()
	//
	// 	h.mutex.Lock()
	// 	h.clients[client] = true
	// 	h.linkToClient(client)
	// 	h.mutex.Unlock()
	//
	// 	// and start a routine to interpret the client's messages
	// 	// NOTE, the push message from server will be sent at hub routine.
	//
	// 	go client.writePump()
	// 	go client.readPump()
	//
	// 	size := int(unsafe.Sizeof(*client))
	// 	logrus.Debugf("=== [ws] new client in. %v clients x %v bytes. %s", h.clients, size, client.userAgent) // and log it for debugging
	// } else {
	// 	h.mutex.RUnlock()
	// 	logrus.Warnf("=== [ws] new client had existed. %v", client)
	// }

	if ok, _ := h.clients[client]; !ok {
		h.clients[client] = true
		logrus.Warnf("=== [ws] new client in. %v", client)
	} else {
		logrus.Warnf("=== [ws] new client had existed. %v", client)
	}
}

func (h *WsHub) onDeregister(client *WsClient) {
	// h.mutex.RLock()
	// if _, ok := h.clients[client]; ok {
	// 	h.mutex.RUnlock()
	//
	// 	h.mutex.Lock()
	// 	delete(h.clients, client)
	// 	h.unlinkClient(client)
	// 	h.mutex.Unlock()
	//
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			err, _ := r.(error)
	// 			logrus.Errorln("Websocket error:", err)
	// 		}
	// 	}()
	//
	// 	client.Close()
	// 	logrus.Printf("=== [ws] the client leaved (%s).", client.userAgent) // and log it for debugging
	// 	client = nil
	// } else {
	// 	h.mutex.RUnlock()
	// 	logrus.Debugf("=== [ws] the client unregistered (%s).", client.userAgent)
	// }

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		client.Close()
		logrus.Printf("=== [ws] the client leaved (%v).", client) // and log it for debugging
		client = nil
	} else {
		logrus.Debugf("=== [ws] the client unregistered (%v).", client)
	}
}
