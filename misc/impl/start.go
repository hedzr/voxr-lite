/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/im_grpc_server"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/internal"
	"github.com/hedzr/voxr-lite/internal/config"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/hedzr/voxr-lite/misc/impl/apps"
	"github.com/hedzr/voxr-lite/misc/impl/filters"
	"github.com/hedzr/voxr-lite/misc/impl/mq"
	"github.com/hedzr/voxr-lite/misc/impl/service"
	"github.com/hedzr/voxr-lite/misc/impl/ws"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var gs *grpc.Server

func InitDaemon(onAppStart func(cmd *cmdr.Command, args []string) (err error), onAppExit func(cmd *cmdr.Command, args []string)) {
	daemon.Enable(
		internal.NewDaemon(doStart, doDeregister, &handlers{}, buildH2Routes),
		modifier, onAppStart, onAppExit)
}

func doStart() {
	logrus.Infof("BEGIN of doStart, port = %v", vxconf.GetIntR("server.port", 2300))
	logrus.Infof("realStart at %v", vxconf.GetIntR("server.port", 2300))

	logrus.Infof("    Starting MQ Manager...")
	mq.Start(internal.AppStopCh())

	logrus.Infof("    Starting JWT Manager...")
	redis_op.JwtInit()

	logrus.Infof("    Starting DB Manager...")
	config.InitDBConn()
	if err := dbe.OpenDbConnection(); err != nil {
		logrus.Errorf("CAN'T OPEN DATABASE: %v", err)
	}

	logrus.Infof("    Starting Apps Manager...")
	apps.Start()

	logrus.Infof("    Starting Filters Manager...")
	filters.Start()

	logrus.Infof("    Starting WS Clients Manager...")
	// run hub service for chatting WebSocket message processing
	// chat.SetMaxMessageSize(vxconf.GetIntR("server.vx-core.websocket.maxMessageSize", 4096))
	ws.StartHub()

	logrus.Infof("    Starting gRPC Manager...")
	scheduler.Start()
	internalMode := vxconf.GetBoolR("internal-mode", false)
	if vxconf.GetStringR("runmode", "devel") != "devel" {
		internalMode = false
	}
	scheduler.AddSynonyms(true, map[string]string{"vx-user": "vx-auth"})
	if internalMode {
		// internal local debug only
		scheduler.AddSynonyms(false, map[string]string{
			"vx-user":    "vx-misc",
			"vx-auth":    "vx-misc",
			"vx-storage": "vx-misc",
		})
	}

	logrus.Infof("    Starting Main gRPC Service...")
	gs = im_grpc_server.StartServer(func(server *grpc.Server) {
		x := service.NewAuthServerV11()
		v10.RegisterUserActionServer(server, x)
		v10.RegisterFriendActionServer(server, x)
		v10.RegisterUserContactServer(server, x)

		v10.RegisterImOrgServer(server, service.NewImOrgService())
		v10.RegisterImTopicServer(server, service.NewImTopicService())
		v10.RegisterImMsgServer(server, service.NewImMsgService())

		// v10.RegisterUserActionServer(server, &service.AuthServer{})
		// v10.RegisterFriendActionServer(server, &service.FriendServer{})
		// v10.RegisterUserContactServer(server, &service.AuthServerV11{})
	})

	// logrus.Infof("    Starting XS Server...")
	// // starting RESTful API service, with WebSocket context handler
	// e := xs.New(restful.New())
	// // xs.SetRegistrarOkCallback(scheduler.RegistrarOkCallback)
	// // xs.SetRegistrarChangesHandler(scheduler.RegistrarChangesHandler)
	// // enter the XsServer main loop
	// defer e.Start()() // will block here //fmt.Println(" END")

	logrus.Infof("END OF doStart.")
}

func doDeregister() {
	// if gs != nil {
	// 	gs.GracefulStop()
	// }

	logrus.Infof("    Stopping WS Clients Manager...")
	ws.StopHub()

	logrus.Infof("    Stopping gRPC services Manager...")
	im_grpc_server.StopServer(gs)

	logrus.Infof("    Stopping gRPC services Manager...")
	scheduler.Stop()

	logrus.Infof("    Stopping Apps Manager...")
	apps.Stop()

	logrus.Infof("    Stopping Filters Manager...")
	filters.Stop()

	logrus.Infof("    Stopping DB Manager...")
	dbe.CloseDbConnection()

	logrus.Infof("    Stopping MQ Manager...")
	mq.Stop()

	// logrus.Infof("    Stop gRPC Client...")
	// coco.StopClient()
	// logrus.Infof("    Stop gRPC Server...")
	// grpc.StopServer()
	// logrus.Infof("    Stop CHAT Service...")
	// chat.StopHub()
	//
	// redis_op.HashDel(fmt.Sprintf("%s:instances", cli_common.AppName), id.GetInstanceId())
	//
	// // redis cluster
	// redis_op.Stop()

	// xs.Deregister() // deregister self from consul/etcd registrar (13.registrar.yml)
}
