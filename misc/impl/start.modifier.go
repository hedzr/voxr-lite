/*
 * Copyright © 2019 Hedzr Yeh.
 */

package impl

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func modifier(daemonServerCommands *cmdr.Command) *cmdr.Command {
	if startCmd := cmdr.FindSubCommand("start", daemonServerCommands); startCmd != nil {
		startCmd.PreAction = onServerPreStart
		startCmd.PostAction = onServerPostStop
	}

	return daemonServerCommands
}

func onServerPostStop(cmd *cmdr.Command, args []string) {
	logrus.Debug("onServerPostStop")
}

// onServerPreStart is earlier than onAppStart.
func onServerPreStart(cmd *cmdr.Command, args []string) (err error) {
	earlierInitLogger()
	logrus.Debug("onServerPreStart")
	return
}

func earlierInitLogger() {
	l := "OFF"
	if !vxconf.IsProd() {
		l = "DEBUG"
	}
	l = vxconf.GetStringR("server.logger.level", l)
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
