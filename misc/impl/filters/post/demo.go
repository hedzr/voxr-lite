/*
 * Copyright © 2019 Hedzr Yeh.
 */

package main

import (
	"encoding/json"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-lite/misc/impl/filters"
	"github.com/sirupsen/logrus"
)

/*

   go build -buildmode=plugin -o demo-pre-filter.so demo.go
   go build -buildmode=plugin -o demo-post-filter.so demo.go

*/

type vxPlugin struct{}

var VxPlugin = vxPlugin{}
var config = &filters.Config{
	Name:        "default-post-filter",
	Icon:        "",
	Version:     "0.0.1",
	Author:      "LAOYE",
	Copyright:   "",
	CaredEvents: "EvMsgPulling", // EvMsgPulling
	Mode:        1,
	Permissions: "",
	Tags:        "",
	Keyword:     "",
	Website:     "",
	Logo:        "",
	Info:        "",
	Cover:       "",
	HelpPage:    "",
	TermsPage:   "",
	Privacy:     "",
}

// func init() {
// 	var x apps.VxPlug = VxPlugin
// 	// x, ok := VxPlug.(apps.VxPlug)
// 	fmt.Println("demo-plugin", x)
// 	// x.OnLoad()
// }

func (vx vxPlugin) Config() *filters.Config {
	return config
}

// post-filter 的首要用途在于屏蔽敏感字。
//
// 通过和敏感字表格进行比对，并自动 mask 这些字样，post-filters 在IM向终端用户提供消息内容列表之前，预先屏蔽敏感字样。
//
// 相似的，可以进一步地开发其他用途的显示前预处理过滤器。
//
func (vx vxPlugin) OnCall(p filters.Plugin, args *filters.Args) (ret *models.Msg, err error) {
	logrus.WithFields(logrus.Fields{
		"Event":  args.Ge,
		"ts":     args.Timestamp,
		"msg":    args.Msg,
		"Filter": p}).
		Infof("[demo-post-filter] got oncall")
	ret = args.Msg
	return
}

func (vx vxPlugin) OnEvent(p filters.Plugin, args *filters.Args) (ret *models.Msg, err error) {
	var b []byte
	b, err = json.Marshal(args.Msg)
	logrus.WithFields(logrus.Fields{
		"Event":       args.Ge,
		"ts":          args.Timestamp,
		"consumerTag": args.ConsumerTag,
		"msg":         args.Msg,
		"body":        b,
		"Filter":      p}).
		Infof("[demo-post-filter] got onevent")
	ret = args.Msg
	return
}

func (vx vxPlugin) OnUnload() (err error) {
	logrus.Infof("test for demo post plugin end.")
	return
}

func (vx vxPlugin) OnLoad() (err error) {
	logrus.Infof("test for demo post plugin loaded. vx=%v", vx)
	return
}
