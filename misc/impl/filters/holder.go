/*
 * Copyright © 2019 Hedzr Yeh.
 */

package filters

import (
	"errors"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-common/dc"
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
		preFilters     map[string]Plugin
		postFilters    map[string]Plugin
		exited         bool
		exitCh         chan bool
		externalExitCh chan bool
		appAdded       chan *models.Filter
		appRemoved     chan *models.Filter
		appUpdated     chan *models.Filter
	}

	Config struct {
		Name        string
		Icon        string
		Version     string
		Author      string
		Copyright   string
		CaredEvents string
		Permissions string
		Tags        string //
		Keyword     string //
		Mode        int    // 0: Pre-filter, 1: Post filter
		Website     string // url
		Logo        string // url
		Info        string // short detail about this filter
		Cover       string // url for cover image
		HelpPage    string // url
		TermsPage   string // url
		Privacy     string // url
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
		make(map[string]Plugin),
		true,
		make(chan bool, 3),
		nil,
		make(chan *models.Filter, 10),
		make(chan *models.Filter, 10),
		make(chan *models.Filter, 10),
	}
}

func (h *holder) loader() {
	// load all filters
	h.loadFilters()
	// starting the run looper
	go h.run()
	// monitor apps add/remove global events
	h.monitorEvents()
}

func (h *holder) loadFilters() {
	dx := dao.NewFilterDao()

	if ret, err := dx.ListFast("1=1"); err != nil {
		logrus.Fatalf("[filters] CAN'T load apps from DB: %v", err)
	} else {
		cnt := 0
		for _, r := range ret {
			if len(r.Name) == 0 || len(r.Callback) == 0 {
				continue
			}

			if strings.HasPrefix(r.Callback, "file://") {
				if h.loadPlugin(r, r.Callback[7:]) {
					cnt++
				}
			}
		}
		logrus.Debugf("[filters] %v filters loaded.", cnt)
	}
}

func (h *holder) loadPlugin(r *models.Filter, file string) (ok bool) {
	if tool.FileExists(file) {
		if p, err := plugin.Open(file); err != nil {
			logrus.WithFields(logrus.Fields{"Err": err}).Warnf("[apps] CAN'T load filters' plugin '%v'", r.Name)

		} else {
			if sym, err := p.Lookup("VxPlugin"); err == nil {
				logrus.Debugf("[filters] 'VxPlugin' is: %v", sym)

				if entry, ok := sym.(VxPlug); ok {
					_ = entry.OnLoad()
					// cfg := entry.Config()
				}

				if r.Mode == 0 {
					h.preFilters[r.Name] = NewFilterPlugin(r, p, sym)
				} else if r.Mode == 1 {
					h.postFilters[r.Name] = NewFilterPlugin(r, p, sym)
				} else {
					h.preFilters[r.Name] = NewFilterPlugin(r, p, sym)
					h.postFilters[r.Name] = NewFilterPlugin(r, p, sym)
				}

				return true
			} else {
				logrus.Warnf("[filters] CANT load app plugin 'VxPlugin' symbol '%v': %v", r.Name, err, errors.New("x"))
			}
		}
	}
	return
}

func (h *holder) eventsHandlerForFilters(d amqp.Delivery) {
	// key := d.RoutingKey
	// if strings.HasPrefix(key, "fx.im.ev.") {
	// 	key = key[9:]
	// }
	// keyInt := util.GlobalEventNameToInt(key)
	// ge := v10.GlobalEvents(keyInt)
	// logrus.Debugf(" [x][filters] %v (%v), ge: %v", d.RoutingKey, keyInt, ge)
	//
	// if ge == v10.GlobalEvents_EvMsgIncoming {
	// 	for _, v := range h.preFilters {
	// 		// logrus.Debugf("v: %v", v)
	// 		entry, ok := v.PluginMainEntry().(VxPlug)
	// 		if ok && v.IsCared(ge) {
	// 			// logrus.Debugf("ge hit: %v | plugin: %v", ge, v.Model.Name)
	// 			go func() {
	// 				logrus.Debugf("run plugin: %v, entry=%v", v.Model().Name, entry)
	// 				if _, err := entry.OnEvent(v, &Args{ge, d.Timestamp, d.ConsumerTag, d.Body}); err != nil {
	// 					logrus.Warnf("[x][filters] invoke app '%v' return failed: %v", v.Model().Name, err)
	// 				}
	// 			}()
	// 		}
	// 	}
	// } else if ge == v10.GlobalEvents_EvMsgRead {
	// 	for _, v := range h.postFilters {
	// 		// logrus.Debugf("v: %v", v)
	// 		entry, ok := v.PluginMainEntry().(VxPlug)
	// 		if ok && v.IsCared(ge) {
	// 			// logrus.Debugf("ge hit: %v | plugin: %v", ge, v.Model.Name)
	// 			go func() {
	// 				logrus.Debugf("run plugin: %v, entry=%v", v.Model().Name, entry)
	// 				if _, err := entry.OnEvent(v, &Args{ge, d.Timestamp, d.ConsumerTag, d.Body}); err != nil {
	// 					logrus.Warnf("[x][filters] invoke app '%v' return failed: %v", v.Model().Name, err)
	// 				}
	// 			}()
	// 		}
	// 	}
	// }
}

func CallPre(ge v10.GlobalEvents, msg *models.Msg) (ret *models.Msg, err error) {
	return holderCore.CallPre(ge, msg)
}

func CallPost(ge v10.GlobalEvents, msg *models.Msg) (ret *models.Msg, err error) {
	return holderCore.CallPost(ge, msg)
}

func (h *holder) CallPre(ge v10.GlobalEvents, msg *models.Msg) (ret *models.Msg, err error) {
	for _, v := range h.preFilters {
		entry, ok := v.PluginMainEntry().(VxPlug)
		if ok && v.IsCared(ge) {
			logrus.Debugf("run pre-filter: %v", v.Model().Name)
			msgCopy := new(models.Msg)
			_ = dc.StandardCopier.Copy(msgCopy, msg)
			ret = msg
			if msgCopy, err = entry.OnCall(v, &Args{ge, time.Now(), "", msgCopy}); err != nil {
				logrus.Warnf("[x][filters] invoke app '%v' return failed: %v", v.Model().Name, err)
			}
		}
	}
	ret = msg
	return
}

func (h *holder) CallPost(ge v10.GlobalEvents, msg *models.Msg) (ret *models.Msg, err error) {
	for _, v := range h.postFilters {
		entry, ok := v.PluginMainEntry().(VxPlug)
		if ok && v.IsCared(ge) {
			logrus.Debugf("run post-filter: %v", v.Model().Name)
			if ret, err = entry.OnCall(v, &Args{ge, time.Now(), "", msg}); err != nil {
				logrus.Warnf("[x][filters] invoke app '%v' return failed: %v", v.Model().Name, err)
			}
		}
	}
	ret = msg
	return
}

func (h *holder) monitorEvents() {
	mq.HandleEvent("filters.mgr", mq.DEFAULT_QUEUE_FOR_FILTERS, mq.DEFAULT_CAST, h.eventsHandlerForFilters)
}

func (h *holder) run() {
	ticker := time.NewTicker(20 * time.Second)
	defer func() {
		ticker.Stop()
		logrus.Debug("--- filters mgr run() stopped.")
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
			logrus.Debugf("--- filter run() looper: %v", tm)

		case c := <-h.appAdded:
			logrus.Debugf("--- filter added: %v", c)
		case c := <-h.appRemoved:
			logrus.Debugf("--- filter removed: %v", c)
		case c := <-h.appUpdated:
			logrus.Debugf("--- filter updated: %v", c)

		}
	}
}
