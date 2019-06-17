/*
 * Copyright © 2019 Hedzr Yeh.
 */

package apps

import (
	"github.com/hedzr/voxr-api/api/v10"
	"time"
)

type (
	Args struct {
		// MsgId       string `json:"mid"`
		Ge          v10.GlobalEvents `json:"ge"`
		Timestamp   time.Time        `json:"ts"`
		ConsumerTag string           `json:"ctag"`
		Body        []byte           `json:"body"`
	}
	// EntryFunc func(p *Plugin, ge v10.GlobalEvents, msgId string, timestamp time.Time, consumerTag string, body []byte) (err error)
	// OnLoadFunc func() error
	// OnUnloadFunc func() error

)
