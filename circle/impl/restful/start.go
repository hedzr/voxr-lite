/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"github.com/hedzr/voxr-api/api/v10"
	redis_op "github.com/hedzr/voxr-common/cache"
	"github.com/hedzr/voxr-common/im_grpc_server"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-lite/circle/impl/service"
	"github.com/hedzr/voxr-lite/internal/config"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var gs *grpc.Server

func doStart() {
	logrus.Infof("BEGIN of doStart, port = %v", vxconf.GetIntR("server.port", 2300))

	logrus.Infof("    Starting JWT Manager...")
	redis_op.JwtInit()

	logrus.Infof("    Starting DB Manager...")
	config.InitDBConn()
	if err := dbe.OpenDbConnection(); err != nil {
		logrus.Errorf("CAN'T OPEN DATABASE: %v", err)
	}

	logrus.Infof("    Starting GRPC Manager...")
	scheduler.Start()
	// internalMode
	internalMode := vxconf.GetBoolR("internal-mode", false)
	if vxconf.GetStringR("runmode", "devel") != "devel" {
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

	logrus.Infof("    Starting Main GRPC Service...")
	gs = im_grpc_server.StartServer(func(server *grpc.Server) {
		x := service.NewImCircleService()
		v10.RegisterImCircleServer(server, x)

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

	logrus.Infof("    Stopping GRPC services Manager...")
	im_grpc_server.StopServer(gs)

	logrus.Infof("    Stopping GRPC services Manager...")
	scheduler.Stop()

	logrus.Infof("    Stopping DB Manager...")
	dbe.CloseDbConnection()

	// xs.Deregister() // deregister self from consul/etcd registrar (13.registrar.yml)
}
