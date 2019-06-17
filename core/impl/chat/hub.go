/*
 * Copyright © 2019 Hedzr Yeh.
 */

package chat

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/core/impl/service"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type (
	Hub struct {
		mutex      *sync.RWMutex
		clients    map[*WsClient]bool
		clientsMap map[string]*WsClient //
		// commands   map[core.MainOp]grpc.CmdFunc
		exitCh     chan bool
		exited     bool
		dumpTicker *time.Ticker
		broadcast  chan []byte
		register   chan *WsClient
		unregister chan *WsClient
		// mesgReading chan *TextMsg
		// mqttReading chan *BinaryMsg
		ppttReading chan *BinaryMsg
		ppttPushing chan *newMsgIncoming
	}

	// the push message structure
	newMsgIncoming struct {
		uid uint64
		nm  *v10.NotifyMessage
	}
)

var hub = Hub{
	mutex:      new(sync.RWMutex),
	clients:    make(map[*WsClient]bool),
	clientsMap: make(map[string]*WsClient),
	exitCh:     make(chan bool),
	exited:     true,
	broadcast:  make(chan []byte),
	register:   make(chan *WsClient),
	unregister: make(chan *WsClient),
	// mesgReading: make(chan *TextMsg),
	// mqttReading: make(chan *BinaryMsg),
	ppttReading: nil,
	ppttPushing: nil,
}

func (h *Hub) OnConfigReloaded() {
	dumpDuration := vxconf.GetDurationR("server.websocket.log.dump-duration", 16*time.Second)
	if dumpDuration < 5*time.Second {
		dumpDuration = 8 * time.Second
	}
	if h.dumpTicker != nil {
		h.dumpTicker.Stop()
	}
	h.dumpTicker = time.NewTicker(dumpDuration)
}

func (h *Hub) run() {
	h.preInitCommands()
	loadMaxMessageSize()

	recvBufferSize := vxconf.GetIntR("server.websocket.pptt.recv-queue.size", 64)
	if recvBufferSize > 4 && recvBufferSize < 32768 {
		h.ppttReading = make(chan *BinaryMsg, recvBufferSize)
	}
	pushBufferSize := vxconf.GetIntR("server.websocket.pptt.push-queue.size", 64)
	if pushBufferSize > 4 && pushBufferSize < 32768 {
		h.ppttPushing = make(chan *newMsgIncoming, pushBufferSize)
	}

	h.OnConfigReloaded()
	cmdr.AddOnConfigLoadedListener(h)
	defer func() {
		cmdr.RemoveOnConfigLoadedListener(h)
		h.dumpTicker.Stop()
		logrus.Debugf("    [ws] chat hub.run() stopped.")
	}()
	// cmd.AddOnConfigReloadedListener(func() {
	// 	dd := loadDumpDuration()
	// 	if dumpDuration != dd {
	// 		ticker.Stop()
	// 		ticker = time.NewTicker(dumpDuration)
	// 	}
	// })
	h.runner()
}

func (h *Hub) runner() {
	defer func() {
		logrus.Debugf("    [ws] chat hub.runner() stopping...")
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorf("    [ws] chat hub.runner() recover err: %v", err)
		}
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
			logrus.Debugf("    [ws] deregistering client: %v", client)
			h.onDeregister(client)

		case binMsg := <-h.ppttReading:
			// logrus.Debugf("ppttReading: %v", *binMsg)
			h.ppttDoProcessMsg(binMsg)

		case msg := <-h.ppttPushing:
			// logrus.Debugf("ppttDoPushing: %v", *msg)
			h.ppttDoPushing(msg)

		case users := <-service.UsersNeedNotified:
			logrus.Debugf("    [ws] UsersNeedNotified: those users need be notified: %v", *users)
			h.onNotifyUsersNewMsgIncoming(users)

		case message := <-h.broadcast:
			h.__b(message)

		case tm := <-h.dumpTicker.C:
			h.__dump(tm)

		}
	}
}

func (h *Hub) __b(message []byte) {
	logrus.Debugf("    [ws] broadcasting text message: %v", string(message))
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for client := range h.clients {
		client._postTxtMsg(message)
	}
}

func (h *Hub) __dump(tm time.Time) {
	if vxconf.GetBoolR("server.websocket.log.dump-clients", false) {
		logrus.Debugf("--- ws clients at %v: ", tm)
		h.mutex.RLock()
		defer h.mutex.RUnlock()
		for client := range h.clients {
			logrus.Debugf("      clients: %v / %v", client.deviceId, client.userAgent)
		}
	} else {
		logrus.Debugf("--- ws clients at %v: %v clients", tm, len(h.clients))
	}
}

func (h *Hub) stop() {
	h.mutex.RLock()
	for k, _ := range h.clients {
		k.conn.Close() // it will break the for loop in ws_hello()
		k.conn = nil
	}
	h.mutex.RUnlock()

	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients = make(map[*WsClient]bool)

	if !h.exited {
		h.exitCh <- true
	}
}

func StartHub() {
	go hub.run()
}

func StopHub() {
	hub.stop()
}
