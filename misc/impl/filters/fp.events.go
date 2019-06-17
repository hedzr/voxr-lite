/*
 * Copyright © 2019 Hedzr Yeh.
 */

package filters

import (
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/util"
	"strings"
)

func (ap *FilterPlugin) IsCared(event v10.GlobalEvents) bool {
	if _, ok := ap.caredEvents[event]; ok {
		return true
	}
	return false
}

func (ap *FilterPlugin) buildEventsMap() {
	ap.caredEvents = make(map[v10.GlobalEvents]bool)
	for _, s := range strings.Split(ap.model.Events, ",") {
		keyInt := util.GlobalEventNameToInt(s)
		if keyInt != v10.GlobalEvents_EvEmpty {
			ap.caredEvents[v10.GlobalEvents(keyInt)] = true
		} else if s == "*" {
			for _, v := range v10.GlobalEvents_value {
				if v != int32(v10.GlobalEvents_EvEmpty) {
					ap.caredEvents[v10.GlobalEvents(keyInt)] = true
				}
			}
		}
	}
}
