/*
 * Copyright © 2019 Hedzr Yeh.
 */

package rand

import (
	"errors"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"math/rand"
)

type (
	RandomBalancer struct {
		// lastPicked int
	}
)

func New() balance2.Balancer {
	return &RandomBalancer{}
}

func (s *RandomBalancer) Name() string {
	return balance2.RANDOM
}

func (s *RandomBalancer) Clone() balance2.Balancer {
	return &RandomBalancer{}
}

func (s *RandomBalancer) UpdateScales(scales map[string]int) (b balance2.Balancer) {
	b = s
	return
}

func (s *RandomBalancer) Pick(peers []*balance2.ServicePeer) (picked *balance2.ServicePeer, err error) {
	if len(peers) == 0 {
		err = errors.New("no instance")
		return
	}

	ix := rand.Intn(len(peers))
	picked = peers[ix]
	return
}
