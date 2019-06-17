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
	Name:        "default-pre-filter",
	Icon:        "",
	Version:     "0.0.1",
	Author:      "LAOYE",
	Copyright:   "",
	CaredEvents: "EvMsgIncoming", // EvMsgPulling
	Permissions: "",
	Tags:        "",
	Keyword:     "",
	Mode:        0,
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

// pre-filter 不应该修改 args.Msg 的内容，它必须原样返回该消息内容。
//
// pre-filter 被用于对特定关键字、标签进行解读和扫描，如有必要，也可以记录 usages 到大数据端点。
// 也可以用 pre-filter 扫描涉黄字样，并统计出重点监控名单。
// 等等。
//
// Filters 并不专门针对某个 org/topic/user 的 Msg，它是面向所有往来的 Msgs。
// 因此 TopicFilter 对象以及相应的表结构都没有被使用。
//
//
func (vx vxPlugin) OnCall(p filters.Plugin, args *filters.Args) (ret *models.Msg, err error) {
	logrus.WithFields(logrus.Fields{
		"Event":  args.Ge,
		"ts":     args.Timestamp,
		"msg":    args.Msg,
		"Filter": p}).
		Infof("[demo-pre-filter] got oncall")
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
		Infof("[demo-pre-filter] got onevent")
	ret = args.Msg
	return
}

func (vx vxPlugin) OnUnload() (err error) {
	logrus.Infof("test for demo pre plugin end.")
	return
}

func (vx vxPlugin) OnLoad() (err error) {
	logrus.Infof("test for demo pre plugin loaded. vx=%v", vx)
	return
}
