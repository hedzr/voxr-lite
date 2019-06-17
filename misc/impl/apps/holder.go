/*
 * Copyright © 2019 Hedzr Yeh.
 */

package apps

import (
	"errors"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-api/util"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-lite/misc/impl/dao"
	"github.com/hedzr/voxr-lite/misc/impl/mq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"plugin"
	"strings"
	"time"
)

type (
	holder struct {
		apps           map[string]Plugin
		exited         bool
		exitCh         chan bool
		externalExitCh chan bool
		appAdded       chan *models.App
		appRemoved     chan *models.App
		appUpdated     chan *models.App
	}

	Plugin interface {
		Model() *models.TopicApp //
		IsPlugin() bool
		Plugin() *plugin.Plugin
		PluginMainEntry() plugin.Symbol
		IsCared(event v10.GlobalEvents) bool
	}
)

var holderCore *holder

func Start() {
	holderCore = newHolder()
	go holderCore.loader()
}

func Stop() {
	if !holderCore.exited {
		holderCore.exitCh <- true
	}
}

func newHolder() *holder {
	return &holder{
		make(map[string]Plugin),
		true,
		make(chan bool, 3),
		nil,
		make(chan *models.App, 10),
		make(chan *models.App, 10),
		make(chan *models.App, 10),
	}
}

func (h *holder) loader() {
	// load all apps
	h.loadApps()
	// starting the run looper
	go h.run()
	// monitor apps add/remove global events
	h.monitorEvents()
}

// TODO 装入所有Topics的关联Apps到内存是不可能的事情，因此需要一个分批装载、或者不装载的方案。
//
// 总的Apps或许是有一个限度的，例如 8000K 个，此时即使仅装载 apps 到内存也是不可能的事情，只能实时检索 TopicApp 表，实时装载对应的 app 并在执行完毕之后卸载该 app。
// 由于 golang 并不支持 unload plugin 以释放相应的资源，因此实时检索、装载、执行、卸载的流程并没有能够复用有限的内存。
// 因此，需要进一步地解决问题：
// 1. 定义 vx-executor 微服务，负责 apps 组件的载入和执行。
// 2. 采用实时检索的策略，在 vx-executor 微服务中载入 apps 并执行。
// 3. 当 vx-executor 装载的 apps 超过一个阈值，例如 8K 个，则通知 管理器停止装载新的 apps，并启动 devops 伸缩机制，建立 vx-executor 的新实例（新的vm-host, 或者新的 container 资源）
// 4. 所以 vx-executor 的集群实质上实现了分批装载全部 apps 的效果。
// 5. 假设每个 vm/container (设为 8G RAM, 8Core) 的 vx-executor 能够容纳和管理 80K 个 apps 实例，那么100台 vm/containers 将能够承载 8000K 个 apps 实例，已经具备足够的可行性。
//
// 当前，仅实现了基本的装载算法，在 apps 不超过 8K 个之前，暂时不考虑实现上述的方案。
//
// 第二，对于提供RESTful回调接口的 app，则考虑不同的策略：
//   此时，3-10 个 vx-executor 实例应该初步满足需要了，这个小型集群负责发出 RESTful 调用、收集结果，实际的 apps 被部署在 IM 核心之外，其消耗的资源不再核心内被考虑。
//   IM平台可以为 apps 提供宿主的功能，用户开发了 app 之后可以采用 IM 运营商所提供的服务器资源完成部署，这时已经去到商业收费的谈话场景，因此这里不再深入探讨了。
//
func (h *holder) loadApps() {
	dx := dao.NewTopicAppDao()

	if ret, err := dx.ListEager("1=1"); err != nil {
		logrus.Fatalf("[apps] CAN'T load apps from DB: %v", err)
		return
	} else {
		cnt := 0
		for _, r := range ret {
			if len(r.App.Name) == 0 || len(r.App.Callback) == 0 {
				continue
			}

			if strings.HasPrefix(r.App.Callback, "file://") {
				if h.loadPlugin(r, r.App.Callback[7:]) {
					cnt++
				}
			} else if strings.HasPrefix(r.App.Callback, "http://") {
				// RESTful API style
				if h.loadRESTfulCB(r, r.App.Callback[7:]) {
					cnt++
				}
			} else if strings.HasPrefix(r.App.Callback, "https://") {
				// RESTful API style
				if h.loadRESTfulCB(r, r.App.Callback[8:]) {
					cnt++
				}
			}
		}
		logrus.Debugf("[apps] %v apps loaded.", cnt)
	}
}

func (h *holder) loadRESTfulCB(r *models.TopicApp, file string) (ok bool) {
	return
}

func (h *holder) loadPlugin(r *models.TopicApp, file string) (ok bool) {
	if tool.FileExists(file) {
		if p, err := plugin.Open(file); err != nil {
			logrus.WithFields(logrus.Fields{"Err": err}).Warnf("[apps] CAN'T load app plugin '%v' from: '%v'", r.App.Name, file)

		} else {
			if sym, err := p.Lookup("VxPlugin"); err == nil {
				logrus.Debugf("[apps] 'VxPlugin' is: %v", sym)
				h.apps[r.App.Name] = NewAppPlugin(r, p, sym)

				if entry, ok := h.apps[r.App.Name].PluginMainEntry().(VxPlug); ok {
					_ = entry.OnLoad()
				}

				return true
			} else {
				logrus.Warnf("[apps] CANT load app plugin 'VxPlugin' symbol '%v': %v", r.App.Name, err, errors.New("x"))
			}
		}
	}
	return
}

func (h *holder) eventsHandlerForApps(d amqp.Delivery) {
	key := d.RoutingKey
	if strings.HasPrefix(key, "fx.im.ev.") {
		key = key[9:]
	}
	keyInt := util.GlobalEventNameToInt(key)
	ge := v10.GlobalEvents(keyInt)
	logrus.Debugf(" [x][apps] %v (%v), ge: %v", d.RoutingKey, keyInt, ge)

	for _, v := range h.apps {
		// logrus.Debugf("v: %v", v)
		entry, ok := v.PluginMainEntry().(VxPlug)
		if ok && v.IsCared(ge) {
			// logrus.Debugf("ge hit: %v | plugin: %v", ge, v.Model.Name)
			go func() {
				logrus.Debugf("run plugin: %v, entry=%v", v.Model().App.Name, entry)
				if err := entry.OnEvent(v, &Args{ge, d.Timestamp, d.ConsumerTag, d.Body}); err != nil {
					logrus.Warnf("[x][apps] invoke app '%v' return failed: %v", v.Model().App.Name, err)
				}
			}()
		}
	}
}

func (h *holder) monitorEvents() {
	mq.HandleEvent("apps.mgr", mq.DEFAULT_QUEUE_FOR_APPS, mq.DEFAULT_CAST, h.eventsHandlerForApps)
}

func (h *holder) run() {
	ticker := time.NewTicker(20 * time.Second)
	defer func() {
		ticker.Stop()
		logrus.Debug("--- apps mgr run() stopped.")
	}()

	for {
		select {
		case e := <-h.exitCh:
			if e {
				Stop()
			}
		case e := <-h.externalExitCh:
			if e {
				return
			}

		case tm := <-ticker.C:
			logrus.Debugf("--- app run() looper: %v", tm)

		case c := <-h.appAdded:
			logrus.Debugf("--- app added: %v", c)
		case c := <-h.appRemoved:
			logrus.Debugf("--- app removed: %v", c)
		case c := <-h.appUpdated:
			logrus.Debugf("--- app updated: %v", c)

		}
	}
}
