/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmd

import (
	"github.com/hedzr/cmdr"
	"github.com/hedzr/logex"
	"github.com/sirupsen/logrus"
)

// Entry is app main entry
func Entry() {

	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// logex.Enable()

	// To disable internal commands and flags, uncomment the following codes
	// cmdr.EnableVersionCommands = false
	// cmdr.EnableVerboseCommands = false
	// cmdr.EnableCmdrCommands = false
	// cmdr.EnableHelpCommands = false
	// cmdr.EnableGenerateCommands = false

	if err := cmdr.Exec(buildRootCmd(),
		cmdr.WithLogex(logrus.DebugLevel),
		cmdr.WithWatchMainConfigFileToo(true),
	); err != nil {
		logrus.Errorf("Error: %v", err)
	}

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
