/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/core/impl/chat"
	"github.com/hedzr/voxr-lite/core/impl/grpc"
	"github.com/hedzr/voxr-lite/core/impl/grpc/coco-client/coco"
	"github.com/hedzr/voxr-lite/internal"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
)

func InitDaemon(onAppStart func(cmd *cmdr.Command, args []string) (err error), onAppExit func(cmd *cmdr.Command, args []string)) {
	daemon.Enable(
		internal.NewDaemon(doStart, doDeregister, &handlers{}, buildH2Routes),
		modifier, onAppStart, onAppExit)
}

func doStart() {
	logrus.Infof("BEGIN of doStart, port = %v", vxconf.GetIntR("server.port", 2300))

	logrus.Infof("    Starting Cache Hub...")
	// redis cluster
	redis_op.Start()
	// redis_op.HashPut(fmt.Sprintf("%s:instances", cli_common.AppName), id.GetInstanceId(), os.Getpid())
	redis_op.AddZone()

	logrus.Infof("    Starting gRPC Server/Clients...")
	// starting grpc service
	grpc.StartServer()

	coco.GrpcStartClient() // a sample grpc client for testing

	logrus.Infof("    Starting Chat Client...")
	// run hub service for chatting WebSocket message processing
	chat.SetMaxMessageSize(vxconf.GetIntR("server.websocket.max-size", 4096))
	chat.StartHub()

	logrus.Infof("    Starting gRPC Manager...")
	// grpc service manager
	scheduler.Start()
	// internalMode
	internalMode := vxconf.GetBoolR("internal-mode", false)
	runMode := vxconf.GetStringR("runmode", "devel")
	if runMode != "devel" {
		internalMode = false
	}
	if internalMode {
		// internal local debug only
		scheduler.AddSynonyms(false, map[string]string{
			"vx-user":    "vx-misc",
			"vx-auth":    "vx-misc",
			"vx-storage": "vx-misc",
		})
	} else {
		scheduler.AddSynonyms(true, map[string]string{"vx-user": "vx-auth"})
	}

	// logrus.Infof("    Starting XS Server...")
	// // starting RESTful API service, with WebSocket context handler
	// e := xs.New(restful.New())
	// xs.SetRegistrarOkCallback(scheduler.RegistrarOkCallback)
	// xs.SetRegistrarChangesHandler(scheduler.RegistrarChangesHandler)
	// // enter the XsServer main loop
	// defer e.Start()() // will block here //fmt.Println(" END")

	logrus.Infof("END OF doStart.")
}

func doDeregister() {
	logrus.Infof("    Stopping Chat Service...")
	chat.StopHub()

	logrus.Infof("    Stopping gRPC services Manager...")
	scheduler.Stop()

	logrus.Infof("    Stopping gRPC Client...")
	coco.GrpcStopClient()
	logrus.Infof("    Stopping gRPC Server...")
	grpc.StopServer()

	logrus.Infof("    Stopping Cache Hub...")
	// redis_op.HashDel(fmt.Sprintf("%s:instances", cli_common.AppName), id.GetInstanceId())
	redis_op.DelZone()
	// redis cluster
	redis_op.Stop()

	logrus.Infof("    Stopping DB Manager...")
	dbe.CloseDbConnection()

	// xs.Deregister() // deregister self from consul/etcd registrar (13.registrar.yml)

	logrus.Infof("    Stopped...")
}
