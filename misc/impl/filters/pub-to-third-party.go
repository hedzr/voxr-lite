/*
 * Copyright © 2019 Hedzr Yeh.
 */

package filters

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"plugin"
	"time"
)

type (
	// VxPlug 一个第三方开发的 VxFilter，是一个具备 filters.VxPlug 接口实现的实体。参见 demo/demo.go 的示意性实现。
	VxPlug interface {
		Config() *Config
		OnCall(p Plugin, args *Args) (ret *models.Msg, err error)
		OnEvent(p Plugin, args *Args) (ret *models.Msg, err error)
		OnLoad() error
		OnUnload() error
	}

	Plugin interface {
		Model() *models.Filter
		Config() *Config
		IsPlugin() bool
		Plugin() *plugin.Plugin
		PluginMainEntry() plugin.Symbol
		IsCared(event v10.GlobalEvents) bool
	}

	Args struct {
		// MsgId       string `json:"mid"`
		Ge          v10.GlobalEvents `json:"ge"`
		Timestamp   time.Time        `json:"ts"`
		ConsumerTag string           `json:"ctag"`
		Msg         *models.Msg      `json:"msg"`
	}
)
