/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"github.com/hedzr/cmdr"
	voxr_lite "github.com/hedzr/voxr-lite"
	"github.com/hedzr/voxr-lite/core/impl/grpc/coco-client/coco"
	"github.com/sirupsen/logrus"
)

func main() {
	// cmd.SetAppName("coco-client", "1.0")
	// cmd.SetRealServerStart(coco.RealMain, coco.Deregister)
	// cmd.SetPrintVersion(coco.PrintVersion)
	// // more_cmds.Enable()
	// cmd.Execute()

	root := cmdr.Root("coco-client", voxr_lite.Version).
		// Header("voxr-lite - An HTTP2 server - no version - hedzr").
		Copyright("reserved", "Hedzr").
		Description("", "").
		Examples("").
		PreAction(onAppStart).
		PostAction(onAppExit).
		Action(RealMain)
	rootCmd := root.RootCommand()

	cmdr.Exec(rootCmd)
}

func onAppStart(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debug("onAppStart")
	return
}

func onAppExit(cmd *cmdr.Command, args []string) {
	logrus.Debug("onAppExit")
}

func RealMain(cmd *cmdr.Command, args []string) (err error) {
	coco.RealMain()
	return
}
