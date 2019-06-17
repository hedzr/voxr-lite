/*
 * Copyright © 2019 Hedzr Yeh.
 */

package mq

import (
	"encoding/json"
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/util"
	"github.com/hedzr/voxr-common/mqe"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
	"time"
)

const (
	DEFAULT_CAST         = mqe.IM_EVENT_CAST     // IM事件总线，广播方式
	DEFAULT_BUS          = mqe.IM_EVENT_BUS      // IM事件总线，非广播方式
	DEFAULT_BUS_WEBHOOKS = mqe.IM_HOOK_EVENT_BUS // IM事件总线，WebHooks专用
	DEFAULT_BUS_APPS     = "im_app_event_bus"    //
	DEFAULT_BUS_USERS    = "im_user_event_bus"   //
	DEFAULT_BUS_MSGS     = "im_msg_event_bus"    //

	DEFAULT_EBEX_CAST         = "fx.ex.event_cast"
	DEFAULT_EBEX              = "fx.ex.event_bus"      // 单消息发送
	DEFAULT_EBQ               = "fx.q.event_bus"       // 单消息发送
	DEFAULT_EBQ_HOOKS         = "fx.q.event_bus.hooks" // 单消息发送
	DEFAULT_EBQ_APPS          = "fx.q.event_bus.apps"  // 单消息发送
	DEFAULT_QUEUE_FOR_APPS    = "fx.q.ebq.apps"        // 广播消息监听
	DEFAULT_QUEUE_FOR_FILTERS = "fx.q.ebq.filters"     // 广播消息监听
	DEFAULT_QUEUE_FOR_HOOKS   = "fx.q.ebq.hooks"       // 广播消息监听

	DEFAULT_ROUTING_KEY = "fx.default.routing"

	MIME_JSON = "application/json"
	MIME_TEXT = "text/plain"

	senderDebug = false
)

var (
	mqEngine      *mqe.MqHub
	defaultClient *mqe.MqClient
)

func FanoutEvent(ev v10.GlobalEvents, msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("json marshal error: %v", err)
		return
	}

	mqEngine.Publish(b, DEFAULT_CAST, ev.String(), MIME_JSON)
}

func FanoutEventBytes(msg []byte) {
	mqEngine.Publish(msg, DEFAULT_CAST, DEFAULT_ROUTING_KEY, MIME_JSON)
}

func Publish(routingKey, contentType string, msg []byte) {
	mqEngine.Publish(msg, DEFAULT_CAST, routingKey, contentType)
}

func PublishText(routingKey string, msg string) {
	mqEngine.Publish([]byte(msg), DEFAULT_CAST, routingKey, MIME_TEXT)
}

func RaiseEvent(ev v10.GlobalEvents, msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("json marshal error: %v", err)
		return
	}

	key := fmt.Sprintf("fx.im.ev.%v", ev.String())
	mqEngine.Publish(b, DEFAULT_CAST, key, MIME_JSON)
}

func Notify(ev v10.GlobalEvents, msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("json marshal error: %v", err)
		return
	}

	key := fmt.Sprintf("fx.im.ev.%v", ev.String())
	mqEngine.Publish(b, DEFAULT_BUS, key, MIME_JSON)
}

func Start(stopCh chan struct{}) {

	if mqEngine != nil {
		return
	}

	mqEngine = mqe.StartPublisherDaemon(stopCh)
	// mqEngine.WithDebug(true, DEFAULT_CAST)

	if senderDebug {
		go func() {
			ticker := time.NewTicker(3 * time.Second)
			defer func() {
				ticker.Stop()
			}()

			for {
				select {
				case tm := <-ticker.C:
					msg := fmt.Sprintf("Hello World! %v", tm)
					// PublishText("fx.im.test", msg)
					RaiseEvent(v10.GlobalEvents_EvOrgAdded, map[string]interface{}{"id": 1, "msg": msg})
				}
			}
		}()
	}

	defaultClient = mqe.NewClient(DEFAULT_CAST, stopCh)

	defaultClient.
		// NewConsumerWithQueueName("apps.mgr", DEFAULT_QUEUE_FOR_APPS, DEFAULT_CAST, eventsHandlerForApps).
		NewConsumerWithQueueName("hooks.mgr", DEFAULT_QUEUE_FOR_HOOKS, DEFAULT_CAST, eventsHandlerForHooks).
		// NewClient(DEFAULT_CAST, common.AppExitCh).
		//	NewConsumer("abc", DEFAULT_CAST, func(d amqp.Delivery) {
		//		logrus.Debugf(" [x] %v | %v", string(d.Body), d.Body)
		//	}).
		// NewConsumerWithQueueName("def", "fx.q.recv3", DEFAULT_CAST, func(d amqp.Delivery) {
		// 	logrus.Debugf(" [-] %v | %v", string(d.Body), d.Body)
		// }).
		// NewConsumerWithQueueName("ghi", "fx.q.recv2", DEFAULT_CAST, func(d amqp.Delivery) {
		// 	logrus.Debugf(" [+] %v | %v", string(d.Body), d.Body)
		// }).
		// NewConsumerWith("abc123", "fx.q.recv1", DEFAULT_CAST, true, false, false, false, nil, func(d amqp.Delivery) {
		// 	logrus.Debugf(" [x] %v | %v", string(d.Body), d.Body)
		// }).
		NewConsumer("abc456", DEFAULT_CAST, func(d amqp.Delivery) {
			logrus.Debugf(" [x][abc456] %v", string(d.Body))
		})

}

func Stop() {

	mqEngine.CloseAll()

}

func HandleEvent(consumerName, queueName, busName string, onRecv func(d amqp.Delivery)) {
	if defaultClient != nil {
		defaultClient.NewConsumerWithQueueName(consumerName, queueName, busName, onRecv)
	}
}

// func eventsHandlerForApps(d amqp.Delivery) {
// 	obj := make(map[string]interface{})
// 	err := json.Unmarshal(d.Body, &obj)
// 	if err != nil {
// 		logrus.Errorf("unmarshal json failed: %v", err)
// 	} else {
// 		key := d.RoutingKey
// 		if strings.HasPrefix(key, "fx.im.ev.") {
// 			key = key[9:]
// 		}
// 		keyInt := util.GlobalEventNameToInt(key)
// 		logrus.Debugf(" [x][apps] %v | %v (%v)", obj, d.RoutingKey, keyInt)
// 	}
// }

func eventsHandlerForHooks(d amqp.Delivery) {
	obj := make(map[string]interface{})
	err := json.Unmarshal(d.Body, &obj)
	if err != nil {
		logrus.Errorf("unmarshal json failed: %v", err)
	} else {
		key := d.RoutingKey
		if strings.HasPrefix(key, "fx.im.ev.") {
			key = key[9:]
		}
		keyInt := util.GlobalEventNameToInt(key)
		logrus.Debugf(" [x][hooks] %v | %v (%v)", obj, d.RoutingKey, keyInt)
	}
}
