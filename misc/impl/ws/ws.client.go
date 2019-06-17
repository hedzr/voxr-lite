/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ws

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// ////////////////////////////////////////////////////////

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	// Maximum message size allowed from peer.
	maxMessageSize int64 = 4096
)

func SetMaxMessageSize(n int) {
	if n > 128 && n < 65536 {
		maxMessageSize = int64(n)
	} else {
		logrus.Warnf("    [ws] incorrect maxMessageSize: %v. it should be in 128..64K.")
	}
}

func loadMaxMessageSize() {
	maxMessageSize1 := vxconf.GetIntR("server.websocket.maxMessageSize", int(maxMessageSize))
	SetMaxMessageSize(maxMessageSize1)
}

type TextMsg struct {
	from *WsClient
	text string
}

type BinaryMsg struct {
	from *WsClient
	body []byte
}

// 每个 Client 代表一个具体的 IM 用户的具体的某一登录设备
type WsClient struct {
	// hub      *WsHub
	conn        *websocket.Conn
	userId      uint64
	deviceId    string
	token       string
	textSend    chan []byte
	ppttSend    chan []byte
	needPushing chan *v10.NotifyMessage
	userAgent   string
	cookies     []*http.Cookie
	exited      bool
	exitCh      chan bool
	seq         uint32
}

func NewWsClient(conn *websocket.Conn, uid uint64, did, token string, userAgent string, cookies []*http.Cookie, registerIt bool) *WsClient {
	var queueSize = vxconf.GetIntR("server.websocket.pptt.send-queue.size", 32)

	client := &WsClient{
		conn:        conn,
		userId:      uid,
		deviceId:    did,
		token:       token,
		textSend:    make(chan []byte, queueSize),
		ppttSend:    make(chan []byte, queueSize),
		userAgent:   userAgent,
		cookies:     cookies,
		exited:      false,
		exitCh:      make(chan bool),
		needPushing: make(chan *v10.NotifyMessage),
	}

	if registerIt {
		// notify hub routine the new client connection object is incoming
		hub.register <- client
	}

	return client
}

func (c *WsClient) Close() {
	if c.exited {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorln("    [ws] Websocket client closing error:", err)
		}
	}()

	c.exited = true
	if c.exitCh != nil {
		c.exitCh <- true
	}
	if c.textSend != nil {
		close(c.textSend)
		c.textSend = nil
	}
	if c.ppttSend != nil {
		close(c.ppttSend)
		c.ppttSend = nil
	}
	if c.needPushing != nil {
		close(c.needPushing)
		c.needPushing = nil
	}
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		if err != nil {
			logrus.Warnf("    [ws] close conn for ws-client failed: %v", err)
		}
	}

	logrus.Debugf("    [ws] CLOSED. %v / %v / %v", c.userId, c.deviceId, c.userAgent)
}

func (c *WsClient) monitor() {
	var exp time.Duration = vxconf.GetDurationR("server.websocket.client-expiration", 10*time.Second)
	ticker := time.NewTicker(exp)
	defer func() {
		logrus.Debugf("    [ws] monitor stopped.")
		ticker.Stop()
		hub.unregister <- c
	}()

	for {
		select {
		case <-ticker.C:
			// service.RefreshUserHash(c.userId, c.deviceId, c.token)
		}
	}
}

func (c *WsClient) readPump() {
	defer func() {
		logrus.Debugf("    [ws] readPump stopped. %v", c)
		hub.unregister <- c
	}()

	logrus.Debugf("    [ws] using ws reading limit to maxMessageSize = %v", maxMessageSize)
	c.conn.SetReadLimit(maxMessageSize)
	// // c.conn.SetCompressionLevel(9)
	// // c.conn.EnableWriteCompression(false)
	// // if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
	// // 	logrus.Warnf("error occurs at ws SetReadDeadline/pongWait: %v", err)
	// // }
	// c.conn.SetPongHandler(func(string) error {
	// 	logrus.Debug("    - websocket.SetPongHandler")
	// 	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
	// 		logrus.Warnf("error occurs at ws SetReadDeadline/pongWait: %v", err)
	// 	}
	// 	return nil
	// })
	// c.conn.SetPingHandler(func(message string) error {
	// 	err := c.conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(pingPeriod))
	// 	if err == websocket.ErrCloseSent {
	// 		return nil
	// 	} else if e, ok := err.(net.Error); ok && e.Temporary() {
	// 		return nil
	// 	}
	// 	return err
	// })

	for {
		type_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Debugf("    [ws] unexpected closed at c.conn.ReadMessage: %v | %v, %v, %v", err, c.userId, c.deviceId, c.userAgent)
				if err = c.conn.Close(); err != nil {
					logrus.Warnf("    [ws] warned at c.conn.ReadMessage (close): %v | %v, %v, %v", err, c.userId, c.deviceId, c.userAgent)
				}
				c.conn = nil
			} else {
				logrus.Warnf("    [ws] warned at c.conn.ReadMessage: %v | %v, %v, %v", err, c.userId, c.deviceId, c.userAgent) // i/o timeout or others
			}
			break
		}

		if type_ == websocket.CloseMessage {
			logrus.Debug("    [ws] NOTICED at c.conn.ReadMessage: is websocket.CloseMessage ")
			break
		}

		// if type_ == websocket.PingMessage {
		// 	if err = c.writeB(websocket.PongMessage, nil); err != nil {
		// 		logrus.Warnf("    [ws] warned at writing PongMessage: %v, IGNORED.", err)
		// 	}
		// 	continue
		// }

		if type_ == websocket.BinaryMessage {
			hub.ppttReading <- &BinaryMsg{c, message}
			continue
		}

		if type_ == websocket.TextMessage {
			msg := string(message)
			if c.onTxtMsg(msg) {
				continue
			}

			// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			// logrus.Debugf("    [ws] broadcasting text message: '%v'", msg)
			hub.broadcast <- message
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logrus.Debugf("    [ws] writePump stopped.")
		ticker.Stop()
		hub.unregister <- c
	}()

	// TODO 进一步避免 writeXXX 的超时影响到整个for loop的运作，提高顺序执行的吞吐量
	for {
		select {
		case message, ok := <-c.ppttSend:
			if ok {
				if err := c.writeB(websocket.BinaryMessage, message); err != nil {
					logrus.Warnf("    [ws] error occurs at ws ppttSend/message: %v", err)
					return
				}
			}

		case message, ok := <-c.textSend:
			if ok {
				if err := c.writeB(websocket.TextMessage, message); err != nil {
					logrus.Warnf("    [ws] error occurs at ws send/message: %v", err)
					return
				} else if string(message) != "pong" {
					logrus.Debugf("    [ws] put text [%d]: '%v'", len(message), string(message))
				}
			}

		case msg, ok := <-c.needPushing:
			if ok {
				if bin, err := proto.Marshal(msg); err != nil {
					logrus.Warnf("    [ws] CAN'T marshal push message: %v | %v", err, msg)
				} else {
					if err := c.writeB(websocket.BinaryMessage, bin); err != nil {
						logrus.Warnf("    [ws] error occurs at ws push/message: %v", err)
						return
					}
				}
			}

		case <-ticker.C:
			if err := c.sendPing(); err != nil {
				logrus.Warnf("    [ws] error occurs at ws send/message: %v", err)
				return
			}

		case <-c.exitCh:
			if err := c.sendCloseMessage(); err != nil {
				logrus.Println("    [ws] write close:", err)
			}
			return
		}
	}
}

func (c *WsClient) sendPing() (err error) {
	err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.PingMessage, []byte{})
}

func (c *WsClient) sendCloseMessage() (err error) {
	if c.conn != nil {
		err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		err = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
	return
}
