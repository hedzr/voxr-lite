/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/voxr-common/vxconf"
	voxr_lite "github.com/hedzr/voxr-lite"
	"github.com/hedzr/voxr-lite/circle"
	"github.com/hedzr/voxr-lite/circle/impl/restful"
	"github.com/hedzr/voxr-lite/cli/cmd"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	cmd.MsEntry(buildRootCmd)
}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	daemon.Enable(restful.NewDaemon(), modifier, onAppStart, onAppExit)

	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {

		_ = os.Setenv("APPNAME", root.AppName)
		logrus.Println("$APPNAME : ", os.Getenv("APPNAME"), os.ExpandEnv("$APPNAME"))

		// // app.server.port
		// if cx := cmdr.FindSubCommand("server", &root.Command); cx != nil {
		//
		// 	// logrus.Debugf("`server` command found")
		//
		// 	if flg := cmdr.FindFlag("port", cx); flg != nil {
		// 		flg.DefaultValue = 2913
		//
		// 	} else {
		// 		opt := cmdr.NewCmdFrom(cx)
		//
		// 		opt.NewFlag(cmdr.OptFlagTypeInt).
		// 			Titles("p", "port").
		// 			Description("the port to listen.", "").
		// 			Group("").
		// 			DefaultValue(2913, "PORT")
		// 	}
		//
		// }
	})

	// root

	root := cmdr.Root(circle.AppName, voxr_lite.Version).
		// Header("voxr-lite - An HTTP2 server - no version - hedzr").
		Copyright(circle.Copyright, "Hedzr").
		Description(circle.Desc, circle.LongDesc).
		Examples(circle.Examples).
		PreAction(onAppStart).
		PostAction(onAppExit)
	rootCmd = root.RootCommand()

	return
}

func modifier(daemonServerCommands *cmdr.Command) *cmdr.Command {
	if startCmd := cmdr.FindSubCommand("start", daemonServerCommands); startCmd != nil {
		startCmd.PreAction = onServerPreStart
		startCmd.PostAction = onServerPostStop
	}

	return daemonServerCommands
}

func onAppStart(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debug("onAppStart")
	return
}

func onAppExit(cmd *cmdr.Command, args []string) {
	logrus.Debug("onAppExit")
}

func onServerPostStop(cmd *cmdr.Command, args []string) {
	logrus.Debug("onServerPostStop")
	// // if gs != nil {
	// // 	gs.GracefulStop()
	// // }
	//
	// logrus.Infof("    Stopping GRPC services Manager...")
	// im_grpc_server.StopServer(gs)
	//
	// logrus.Infof("    Stopping GRPC services Manager...")
	// scheduler.Stop()
	//
	// logrus.Infof("    Stopping DB Manager...")
	// dbe.CloseDbConnection()
	//
	// // xs.Deregister() // deregister self from consul/etcd registrar (13.registrar.yml)
}

// onServerPreStart is earlier than onAppStart.
func onServerPreStart(cmd *cmdr.Command, args []string) (err error) {
	earlierInitLogger()
	logrus.Debug("onServerPreStart")
	// logrus.Infof("    Starting JWT Manager...")
	// redis_op.JwtInit()
	//
	// logrus.Infof("    Starting DB Manager...")
	// config.InitDBConn()
	// if err := dbe.OpenDbConnection(); err != nil {
	// 	logrus.Errorf("CAN'T OPEN DATABASE: %v", err)
	// }
	//
	// logrus.Infof("    Starting GRPC Manager...")
	// scheduler.Start()
	// // internalMode := vxconf.GetBoolR("internal-mode", false)
	// // if vxconf.GetStringR("runmode", "devel") != "devel" {
	// // 	internalMode = false
	// // }
	// scheduler.AddSynonyms(true, map[string]string{"vx-user": "vx-auth"})
	// // if internalMode {
	// // 	// internal local debug only
	// // 	scheduler.AddSynonyms(false, map[string]string{
	// // 		"vx-user":    "vx-circle",
	// // 		"vx-auth":    "vx-circle",
	// // 		"vx-storage": "vx-circle",
	// // 	})
	// // }
	//
	// logrus.Infof("    Starting Main GRPC Service...")
	// gs = im_grpc_server.StartServer(func(server *grpc.Server) {
	// 	x := service.NewImCircleService()
	// 	v10.RegisterImCircleServer(server, x)
	//
	// 	// v10.RegisterUserActionServer(server, &service.AuthServer{})
	// 	// v10.RegisterFriendActionServer(server, &service.FriendServer{})
	// 	// v10.RegisterUserContactServer(server, &service.AuthServerV11{})
	// })
	//
	// logrus.Infof("    Starting XS Server...")
	// // starting RESTful API service, with WebSocket context handler
	// e := xs.New(restful.New())
	// // xs.SetRegistrarOkCallback(scheduler.RegistrarOkCallback)
	// // xs.SetRegistrarChangesHandler(scheduler.RegistrarChangesHandler)
	// // enter the XsServer main loop
	// defer e.Start()() // will block here //fmt.Println(" END")
	//
	// logrus.Infof("ENDING")

	return
}

func earlierInitLogger() {
	l := vxconf.GetStringR("server.logger.level", "OFF")
	logrus.SetLevel(stringToLevel(l))
	if l == "OFF" {
		logrus.SetOutput(ioutil.Discard)
	}
}

func stringToLevel(s string) logrus.Level {
	s = strings.ToUpper(s)
	switch s {
	case "TRACE":
		return logrus.TraceLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.FatalLevel
	}
}
