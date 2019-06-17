/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ws

import (
	"github.com/sirupsen/logrus"
	"unsafe"
)

func (h *WsHub) onRegister(client *WsClient) {
	h.mutex.RLock()
	if ok, _ := h.clients[client]; !ok {
		h.mutex.RUnlock()

		h.mutex.Lock()
		h.clients[client] = true
		h.linkToClient(client)
		h.mutex.Unlock()

		// and start a routine to interpret the client's messages
		// NOTE, the push message from server will be sent at hub routine.

		go client.writePump()
		go client.readPump()
		go client.monitor()

		size := int(unsafe.Sizeof(*client))
		logrus.Debugf("=== [ws] new client in. %v clients x %v bytes. %s", h.clients, size, client.userAgent) // and log it for debugging
	} else {
		h.mutex.RUnlock()
		logrus.Warnf("=== [ws] new client had existed. %v", client)
	}
}

func (h *WsHub) onDeregister(client *WsClient) {
	h.mutex.RLock()
	if _, ok := h.clients[client]; ok {
		h.mutex.RUnlock()

		h.mutex.Lock()
		delete(h.clients, client)
		h.unlinkClient(client)
		h.mutex.Unlock()

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		err, _ := r.(error)
		// 		logrus.Errorln("    [ws] Websocket error:", err)
		// 	}
		// }()

		go func() {
			client.Close()
			logrus.Printf("=== [ws] the client leaved,closed,unregistered (%s).", client.userAgent) // and log it for debugging
			client = nil
		}()
	} else {
		h.mutex.RUnlock()
		logrus.Debugf("=== [ws] the client unregistered (%s).", client.userAgent)
	}
}

func (h *WsHub) linkToClient(c *WsClient) {
	h.clientsMap[c.deviceId] = c
}

func (h *WsHub) unlinkClient(c *WsClient) {
	delete(h.clientsMap, c.deviceId)
}
