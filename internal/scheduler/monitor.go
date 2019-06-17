/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-common/kvs/store"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/vxconf/gwk"
	"github.com/sirupsen/logrus"
	"strings"
)

var RegistrarOkCallback = func(r *gwk.Registrar, serviceName string) {
	logrus.Debugf("registrar connected: %s, %v", serviceName, r)

	// 添加自身的主服务条目
	serviceAddDepEntry(serviceName, ":1")

	// 添加依赖服务的条目
	var out = make(map[string]*DepRecord)
	_ = vxconf.LoadSectionTo("server.deps", out)
	for k := range out {
		serviceAddDepEntry(k, ":1")
	}

	// fObj := vxconf.GetR("server.deps")
	// if fObj != nil {
	// 	b, err := yaml.Marshal(fObj)
	// 	if err == nil {
	// 		var out = make(map[string]*DepRecord)
	// 		err = yaml.Unmarshal(b, &out)
	// 		if out != nil && len(out) > 0 {
	// 			for k, _ := range out {
	// 				serviceAddDepEntry(k, ":1")
	// 			}
	// 		}
	// 	}
	// }
}

var RegistrarChangesHandler = func(evType store.Event_EventType, key []byte, value []byte) {
	k := string(key)
	v := string(value)

	// fo non-e3w mode, we strip the leading slash.
	if len(k) > 0 && k[0] == '/' {
		k = k[1:]
	}

	parts := strings.Split(k, "/")
	if parts[0] != "services" {
		parts = parts[1:]
	}
	if parts[0] == "services" {
		if len(parts) > 4 && parts[4] == api.GRPC {
			onUpdatedAx(evType, parts[1], parts[2], parts[3], v)
		} else if len(parts) > 1 {
			logrus.Debugf(">>> registrar changing: %v, %v, %v", store.EvTypeToString(evType), parts[1], v)
		} else {
			logrus.Debugf(">>> registrar changing: %v, %v, %v", store.EvTypeToString(evType), k, v)
		}
	}
}

func onUpdatedAx(evType store.Event_EventType, serviceName, peers, peer, value string) {
	if c, ok := grpcHub.byIds[serviceName]; ok {
		if ov, ok := lastValues[peer]; ok && ov == value {
			return
		} else if ok {
			if ov != value {
				lastValues[peer] = value
				onServiceUpdatedAx(c, serviceName, peer, ov, value)
			}
		} else {
			lastValues[peer] = value
		}
	} else {
		logrus.Debugf(">> Unknown service '%s'/'%s' FOUND, not deps of mine, ignored.", serviceName, peer)
		// refreshServicesSync()
	}
}

func onServiceUpdatedAx(c *GrpcClient, serviceName, peer, oldValue, value string) {
	logrus.Debugf(">> Service '%s'/'%s' changed, %v => %v.", serviceName, peer, oldValue, value)
	grpcHub.refreshClients <- c.Name
}

// 添加依赖服务的关注条目
// 这些条目的 grpc 值的变更将被关注和保持更新
func serviceAddDepEntry(serviceName string, grpcValue string) {
	lastValues[serviceName] = grpcValue
}

var (
	lastValues = make(map[string]string)
)
