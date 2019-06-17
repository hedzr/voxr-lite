/*
 * Copyright © 2019 Hedzr Yeh.
 */

package apps

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"plugin"
)

type (
	AppPlugin struct {
		model       *models.TopicApp // TopicApp 包含 Org's apps or Topic's apps 关联关系
		plugin      *plugin.Plugin
		mainEntry   plugin.Symbol
		caredEvents map[v10.GlobalEvents]bool
	}
)

func NewAppPlugin(model *models.TopicApp, plugin *plugin.Plugin, mainEntry plugin.Symbol) (r *AppPlugin) {
	r = &AppPlugin{
		model:       model,
		plugin:      plugin,
		mainEntry:   mainEntry,
		caredEvents: make(map[v10.GlobalEvents]bool),
	}
	r.buildEventsMap()
	return
}

func (ap *AppPlugin) Model() *models.TopicApp {
	return ap.model
}

func (ap *AppPlugin) IsPlugin() bool {
	return ap.plugin != nil
}

func (ap *AppPlugin) Plugin() *plugin.Plugin {
	return ap.plugin
}

func (ap *AppPlugin) PluginMainEntry() plugin.Symbol {
	return ap.mainEntry
}
