/*
 * Copyright © 2019 Hedzr Yeh.
 */

package filters

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"plugin"
)

type (
	FilterPlugin struct {
		model       *models.Filter
		plugin      *plugin.Plugin
		mainEntry   plugin.Symbol
		caredEvents map[v10.GlobalEvents]bool
	}
)

func NewFilterPlugin(model *models.Filter, plugin *plugin.Plugin, mainEntry plugin.Symbol) (r *FilterPlugin) {
	r = &FilterPlugin{
		model:       model,
		plugin:      plugin,
		mainEntry:   mainEntry,
		caredEvents: make(map[v10.GlobalEvents]bool),
	}
	r.buildEventsMap()
	return
}

func (ap *FilterPlugin) Config() *Config {
	if entry, ok := ap.mainEntry.(VxPlug); ok {
		return entry.Config()
	}
	return nil
}

func (ap *FilterPlugin) Model() *models.Filter {
	return ap.model
}

func (ap *FilterPlugin) IsPlugin() bool {
	return ap.plugin != nil
}

func (ap *FilterPlugin) Plugin() *plugin.Plugin {
	return ap.plugin
}

func (ap *FilterPlugin) PluginMainEntry() plugin.Symbol {
	return ap.mainEntry
}
