/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ws

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type (
	WsHub struct {
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

var hub = WsHub{
	mutex:       new(sync.RWMutex),
	clients:     make(map[*WsClient]bool),
	clientsMap:  make(map[string]*WsClient),
	exitCh:      make(chan bool),
	exited:      true,
	broadcast:   make(chan []byte),
	register:    make(chan *WsClient),
	unregister:  make(chan *WsClient),
	ppttReading: nil,
	ppttPushing: nil,
	// mesgReading: make(chan *TextMsg),
	// mqttReading: make(chan *BinaryMsg),
}

func (h *WsHub) OnConfigReloaded() {
	dumpDuration := vxconf.GetDurationR("server.websocket.log.dump-duration", 16*time.Second)
	if dumpDuration < 5*time.Second {
		dumpDuration = 8 * time.Second
	}
	if h.dumpTicker != nil {
		h.dumpTicker.Stop()
	}
	h.dumpTicker = time.NewTicker(dumpDuration)
}

func (h *WsHub) preInitCommands() {
}

func (h *WsHub) run() {
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

func (h *WsHub) runner() {
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
			logrus.Debugf("    [ws] ppttReading: %v", *binMsg)
			h.ppttDoProcessMsg(binMsg)

		case msg := <-h.ppttPushing:
			logrus.Debugf("    [ws] ppttDoPushing: %v", *msg)
			h.ppttDoPushing(msg)

		// case users := <-service.UsersNeedNotified:
		// 	logrus.Debugf("    [ws] UsersNeedNotified: those users need be notified: %v", *users)
		// 	h.onNotifyUsersNewMsgIncoming(users)

		case msg := <-h.broadcast:
			logrus.Debugf("    [ws] broadcast: %v", string(msg))
			h.__b(msg)

		case tm := <-h.dumpTicker.C:
			h.__dump(tm)

		}
	}
}

func (h *WsHub) ppttDoProcessMsg(bin *BinaryMsg) {
	// // logrus.Debugf("pptt: (%v,%v) -> %v", bin.from.userId, bin.from.deviceId, bin.body)
	// if bin.body[0] == 0xa5 {
	// 	// fn(bin.from, bin.body[1:])
	//
	// 	var lead = int(bin.body[1])
	// 	// var op = v10.Op(int(bin.body[2])) // 预取一字节以判断 MainOp，注意现在约定 MainOp 编号值不大于 128，因此可以以一个字节来完成检测
	// 	var op int64
	// 	var ate int
	// 	op, ate = api.DecodeZigZagInt(bin.body[1:])
	//
	// 	if lead == 8 { // pb tag = 1
	// 		if err := grpc.Instance.PooledInvoke(v10.Op(op), bin.from, bin.body[1:]); err != nil {
	// 			h.writeBack(bin.from, fmt.Sprintf("    [ws] [WARN] Unknown Data Diagram. (op=%v,ate=%v). %v", op,ate, err))
	// 		}
	// 	} else {
	// 		h.writeBack(bin.from, "    [ws] [WARN] Unsupport Data Diagram.")
	// 	}
	// }
}

// // after
// func (h *WsHub) onNotifyUsersNewMsgIncoming(users *service.NotifiedUsers) {
// 	// for _, uid := range users.Users {
// 	// 	push := &v10.TalkingPush{
// 	// 		SubscribeId:         users.ConversationId,
// 	// 		ConversationSection: users.ConversationSection,
// 	// 		SortNum:             users.SortNum,
// 	// 		Ts:                  users.TS, MsgIds: []uint64{users.MsgId,},
// 	// 	}
// 	// 	nm := &v10.NotifyMessage{ProtoOp: v10.Op_NotifyAck, Oneof: &v10.NotifyMessage_Talking{Talking: push,},}
// 	// 	h.ppttPushing <- &newMsgIncoming{uid, nm,}
// 	// }
// }

func (h *WsHub) ppttDoPushing(msg *newMsgIncoming) {
	// find all clients in all zone by uid
	// and push msg.nm to him.
}

func (h *WsHub) writeBack(from *WsClient, msg string) {
	from._writeBack(msg)
	// _ = from.conn.WriteMessage(1, []byte(msg))
}

func (h *WsHub) handshake(from *WsClient, body []byte) {
	//
}

func (h *WsHub) __b(message []byte) {
	logrus.Debugf("    [ws] broadcasting text message: %v", string(message))
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for client := range h.clients {
		client._postTxtMsg(message)
	}
}

func (h *WsHub) __dump(tm time.Time) {
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

func (h *WsHub) stop() {
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
