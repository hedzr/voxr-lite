/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	WsClient struct {
		IsReady           bool
		ReconnectDuration time.Duration
		needReconnect     chan bool
		conn              *websocket.Conn
		token             string
		did               string
		textSend          chan []byte
		ppttSend          chan []byte
		exitCh            chan bool
		uid               uint64
		// interrupt         chan os.Signal
	}
)

// var wsClient *WsClient

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

func loadMaxMessageSize() {
	maxMessageSize1 := vxconf.GetIntR("server.websocket.maxMessageSize", int(maxMessageSize))
	SetMaxMessageSize(maxMessageSize1)
}

func SetMaxMessageSize(n int) {
	if n > 128 && n < 65536 {
		maxMessageSize = int64(n)
	} else {
		logrus.Warnf("incorrect maxMessageSize: %v. it should be in 128..64K.")
	}
}

func (c *WsClient) Init() error {
	// var interrupt chan os.Signal
	// interrupt = make(chan os.Signal, 2)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	c.exitCh = make(chan bool)

	token, did, uid, err := doLogin()
	logrus.Printf("device id: %v\n", did)
	logrus.Printf("token    : %v\n", token)
	if err != nil {
		logrus.Printf("error    : %v\n", err)
		return err
	}
	c.token = token
	c.did = did
	c.uid = uid

	c.ppttSend = make(chan []byte)
	c.textSend = make(chan []byte)
	c.needReconnect = make(chan bool)

	go c.scheduler()
	return nil
}

func (c *WsClient) Reconnect() {
	if c.ppttSend == nil {
		c.ppttSend = make(chan []byte)
	}
	if c.textSend == nil {
		c.textSend = make(chan []byte)
	}
	if c.needReconnect == nil {
		c.needReconnect = make(chan bool)
	}
	time.Sleep(c.ReconnectDuration)
	c.needReconnect <- true
}

func (c *WsClient) CloseFully() {
	c.CloseNoReconnect()
}

func (c *WsClient) CloseNoReconnect() {
	if c.textSend != nil {
		close(c.textSend)
		c.textSend = nil
	}
	if c.ppttSend != nil {
		close(c.ppttSend)
		c.ppttSend = nil
	}
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		if err != nil {
			logrus.Warnf("close conn for ws-client failed: %v", err)
		}
	}
}

func (c *WsClient) Close() {
	c.CloseNoReconnect()
	c.Reconnect()
}

func (c *WsClient) scheduler() {
	defer func() {
		logrus.Debugf("    [ws] scheduler stopped. %v", c)
	}()
	for {
		select {
		case _ = <-c.needReconnect:
			c.Open()
		case _ = <-c.exitCh:
			return
		}
	}
}

func (c *WsClient) Open() {
	u := url.URL{Scheme: "ws", Host: "localhost:7111", Path: "/v1/ws"}
	logrus.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", c.token)},
		"X-Device-Id":   []string{c.did},
	})
	if err != nil {
		logrus.Error("dial:", err)
		time.Sleep(c.ReconnectDuration * 1)
		logrus.Debugf("after open failed, reconnecting...")
		go c.Reconnect()
		return
	}
	reconnectAgain := false
	defer func() {
		c.IsReady = false
		if r := recover(); r != nil {
			err, _ := r.(error)
			logrus.Errorln("Websocket error:", err)

			if reconnectAgain {
				go c.Reconnect()
			}
		}
	}()
	// defer conn.Close()
	c.conn = conn

	done := make(chan struct{})

	c.IsReady = true

	go c.wsWritePump(done)
	go c.wsReadPump(done)

	wshub.register <- c

	// err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	// err = conn.WriteMessage(websocket.TextMessage, []byte("PING"))
	// if err != nil {
	// 	logrus.Fatal("dial:", err)
	// }
	// logrus.Info("text 'ping' sent.")
	// // <-interrupt
}

func (c *WsClient) sendTxtMsg(message []byte) {
	c.textSend <- message
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

func (c *WsClient) sendTextMessage(msg string) (err error) {
	err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err = c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

func (c *WsClient) sendPpttMessage(msg []byte) (err error) {
	var bin []byte // = make([]byte, len(msg)+1)
	bin = append([]byte{0xa5}, msg...)
	err = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err = c.conn.WriteMessage(websocket.BinaryMessage, bin)
	return
}

func (c *WsClient) writeN(ok bool, type_ int, message []byte) (err error) {
	if err = c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}
	if !ok {
		// The hub closed the channel.
		if err = c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			err = nil
		}
		return
	}

	var w io.WriteCloser
	w, err = c.conn.NextWriter(type_)
	if err != nil {
		return
	}
	if _, err = w.Write(message); err != nil {
		return
	}

	// Add queued chat messages to the current websocket message.
	n := len(c.textSend)
	for i := 0; i < n; i++ {
		_, err = w.Write(newline)
		if err == nil {
			_, err = w.Write(<-c.textSend)
		}
		if err != nil {
			return
		}
	}

	err = w.Close()
	return
}

func (c *WsClient) writeB(type_ int, data []byte) (err error) {
	if err = c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		logrus.Warnf("error occurs at ws SetWriteDeadline/writeWait: %v", err)
		return
	}

	if err = c.conn.WriteMessage(type_, data); err != nil {
		logrus.Warnf("error occurs at ws WriteMessage/%d: %v", type_, err)
	}
	return
}

func (c *WsClient) wsReadPump(done chan struct{}) {
	defer func() {
		logrus.Debugf("    [ws] wsReadPump stopped. %v", c)
		close(done)
	}()

	logrus.Debugf("    [ws] using ws reading limit to maxMessageSize = %v", maxMessageSize)
	c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetCompressionLevel(9)
	// c.conn.EnableWriteCompression(false)
	// if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
	// 	logrus.Warnf("error occurs at ws SetReadDeadline/pongWait: %v", err)
	// }
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
		start := time.Now()
		type_, message, err := c.conn.ReadMessage()
		elapsed := time.Since(start)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Debugf("    [ws] [after %v] unexpected closed at c.conn.ReadMessage: %v", elapsed, err)
				if err = c.conn.Close(); err != nil {
					logrus.Warnf("    [ws] warned at c.conn.ReadMessage (close): %v", err)
				}
				c.conn = nil
			} else {
				logrus.Warnf("    [ws] [after %v] warned at c.conn.ReadMessage: %v", elapsed, err) // i/o timeout or others
			}
			break
		}

		if type_ == websocket.CloseMessage {
			logrus.Debug("    [ws] NOTICED at c.conn.ReadMessage: is websocket.CloseMessage ")
			break
		}

		// if type_ == websocket.PingMessage {
		// 	logrus.Debug("    - websocket.PingMessage")
		// 	if err = c.writeB(websocket.PongMessage, nil); err != nil {
		// 		logrus.Warnf("    [ws] warned at writing PongMessage: %v, IGNORED.", err)
		// 	}
		// 	continue
		// }

		// if err != nil {
		// 	logrus.Println("read:", err)
		// 	return
		// }

		if type_ == websocket.TextMessage {
			logrus.Printf("    [ws] recv %v bytes: %s", len(message), string(message))
		} else if type_ == websocket.BinaryMessage {
			if err := c.ppttProcess(message); err != nil {
				logrus.Printf("    [ws] recv %v bytes: %v", len(message), message)
				logrus.Warnf("    [ws] ERR: %v", err)
			}
		}
	}
}

func (c *WsClient) wsWritePump(done chan struct{}) {
	ticker := time.NewTicker(40 * time.Second)
	defer func() {
		logrus.Debugf("    [ws] wsWritePump stopped.")
		ticker.Stop()
		wshub.unregister <- c

		// c.Close() // close websocket connection
		//
		// logrus.Debug("preparing reconnect...")
		// // if c.reconnectAgain {
		// time.Sleep(c.ReconnectDuration)
		// logrus.Debug("reconnecting...")
		// c.Open()
		// // }
	}()

	var err error

	for {
		select {
		case t := <-ticker.C:
			if vxconf.GetBoolR("server.websocket.send.u-text", false) {
				err = c.sendTextMessage(fmt.Sprintf("u:%v", t.String()))
				if err != nil {
					logrus.Errorf("    [ws] ticker sendTextMessage: %v", err)
					return
				}
			}

		case x := <-c.textSend:
			if err = c.sendTextMessage(string(x)); err != nil {
				logrus.Errorf("    [ws] textSend sendTextMessage: %v", err)
				return
			}

		case <-done:
			logrus.Debugf("    [ws] `done` signal triggered. exiting...")
			return

		case <-c.exitCh:
			logrus.Println("    [ws] `exitCh` interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err = c.sendCloseMessage()
			if err != nil {
				logrus.Println("    [ws] write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("    [ws] %s took %v\n", what, time.Since(start))
	}
}

func doTiming() {
	defer elapsed("page")()
	time.Sleep(time.Second * 2)
}
