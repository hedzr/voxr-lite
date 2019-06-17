/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"encoding/json"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/misc/impl/apps"
	"github.com/sirupsen/logrus"
)

/*

   go build -buildmode=plugin -o demo-app-plugin.so demo.go

*/

type vxPlugin struct{}

var VxPlugin = vxPlugin{}

// func init() {
// 	var x apps.VxPlug = VxPlugin
// 	// x, ok := VxPlug.(apps.VxPlug)
// 	fmt.Println("demo-plugin", x)
// 	// x.OnLoad()
// }

func (vx vxPlugin) OnEvent(p apps.Plugin, args *apps.Args) (err error) {
	var obj interface{}
	if args.Ge == v10.GlobalEvents_EvOrgAdded {
		obj = new(models.Organization)
		err = json.Unmarshal(args.Body, &obj)
	} else {
		obj = make(map[string]interface{})
		err = json.Unmarshal(args.Body, &obj)
	}
	logrus.WithFields(logrus.Fields{
		"Event":       args.Ge,
		"ts":          args.Timestamp,
		"consumerTag": args.ConsumerTag,
		"body":        obj,
		"App":         p}).
		Infof("[demo-app] got")
	return
}

func (vx vxPlugin) OnUnload() (err error) {
	logrus.Infof("test for demo app plugin end.")
	return
}

func (vx vxPlugin) OnLoad() (err error) {
	logrus.Infof("test for demo app plugin loaded. vx=%v, fxp=%v", vx, VxPlugin)
	return
}
