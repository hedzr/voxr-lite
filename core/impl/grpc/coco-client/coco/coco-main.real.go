/*
 * Copyright © 2019 Hedzr Yeh.
 */

package coco

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"time"
)

func PrintVersion() {
	fmt.Printf(`%15s Version: %s
	        Githash: %s
	       Build at: %s
`,
		conf.AppName, conf.Version, conf.Githash, conf.Buildstamp)
}

func Deregister() {
	// if wsClient != nil {
	// 	wsClient.Close()
	// 	wsClient = nil
	// }

	// GrpcStopClient()

	WsStop()
}

func RealMain() {
	listen, id, disabled, port := tool.LoadGRPCListen("server.deps.vx-core")
	logrus.Infof("listen, id, disabled, port: %v, %v, %v, %v", listen, id, disabled, port)

	// GrpcStartClient()

	WsStart()
	for i := 1; i <= vxconf.GetIntR("server.websocket.count", 8); i++ {
		WsClientAddNew()
		time.Sleep(3 * time.Second)
	}

	<-daemon.QuitSignals()

	// wsClient = &WsClient{
	// 	ReconnectDuration: 3 * time.Second,
	// }
	// wsClient.Init()
	//
	// // hold here
	// wsClient.Open()

	// shutdown...

	// scheduler.Stop()
	//
	// StopClient()
}
