/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/logex"
	voxr_lite "github.com/hedzr/voxr-lite"
	"github.com/hedzr/voxr-lite/core/impl"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	appName   = "vx-core"
	copyright = "vx-core is a part of voxr-lite microservices."
	longDesc  = "vx-core is a part of voxr-lite microservices."
	examples  = `
$ {{.AppName}} gen shell [--bash|--zsh|--auto]
  generate bash/shell completion scripts
$ {{.AppName}} gen man
  generate linux man page 1
$ {{.AppName}} --help
  show help screen.
`
	overview = ``
)

func main() {
	MsEntry(buildRootCmd)
}

func MsEntry(buildRootCmd func() *cmdr.RootCommand) {
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logex.Enable()

	// To disable internal commands and flags, uncomment the following codes
	// cmdr.EnableVersionCommands = false
	// cmdr.EnableVerboseCommands = false
	// cmdr.EnableCmdrCommands = false
	// cmdr.EnableHelpCommands = false
	// cmdr.EnableGenerateCommands = false

	if err := cmdr.Exec(buildRootCmd()); err != nil {
		logrus.Errorf("Error: %v", err)
	}
}

func buildRootCmd() (rootCmd *cmdr.RootCommand) {
	impl.InitDaemon(onAppStart, onAppExit)

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

	root := cmdr.Root(appName, voxr_lite.Version).
		// Header("voxr-lite - An HTTP2 server - no version - hedzr").
		Copyright(copyright, "Hedzr").
		Description(longDesc, longDesc).
		Examples(examples).
		PreAction(onAppStart).
		PostAction(onAppExit)
	rootCmd = root.RootCommand()

	return
}

func onAppStart(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debug("onAppStart")
	return
}

func onAppExit(cmd *cmdr.Command, args []string) {
	logrus.Debug("onAppExit")
}
