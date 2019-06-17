/*
 * Copyright © 2019 Hedzr Yeh.
 */

package balancers

import (
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/hash"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/rand"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/round-robin"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/version-scale"
	"github.com/sirupsen/logrus"
)

var (
	balancers balance2.Balancers
)

func init() {
	balancers = make(balance2.Balancers)
	Add(rand.New())
	Add(round_robin.New())
	Add(hash.New())
	Add(version_scale.New(nil))
}

func Add(balancer balance2.Balancer) {
	balancers[balancer.Name()] = balancer
}

func UpdateScales(scales map[string]int) {
	for _, v := range balancers {
		v.UpdateScales(scales)
	}
}

func NewWith(type_, subType_ string, scales map[string]int) (b balance2.Balancer) {
	if v, ok := balancers[type_]; ok {
		b = v.Clone()
		b.UpdateScales(scales)
	} else {
		logrus.Fatalf("error balancer name: %v, %v, %v", type_, subType_, scales)
	}
	return
}

func PickPeer(method string, peers []*balance2.ServicePeer) (picked *balance2.ServicePeer, err error) {
	if b, ok := balancers[method]; ok {
		picked, err = b.Pick(peers)
		logrus.Debugf("%v : %v", b, ok)
	}
	return
}
