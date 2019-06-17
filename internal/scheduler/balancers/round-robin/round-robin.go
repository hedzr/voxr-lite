/*
 * Copyright © 2019 Hedzr Yeh.
 */

package round_robin

import (
	"errors"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
)

type (
	RoundRobinBalancer struct {
		currIndex int
	}
)

func New() balance2.Balancer {
	return &RoundRobinBalancer{0}
}

func (s *RoundRobinBalancer) Name() string {
	return balance2.ROUNDROBIN
}

func (s *RoundRobinBalancer) Clone() balance2.Balancer {
	return &RoundRobinBalancer{s.currIndex}
}

func (s *RoundRobinBalancer) UpdateScales(scales map[string]int) (b balance2.Balancer) {
	b = s
	return
}

func (s *RoundRobinBalancer) Pick(peers []*balance2.ServicePeer) (picked *balance2.ServicePeer, err error) {
	if len(peers) == 0 {
		err = errors.New("no instance")
		return
	}

	lens := len(peers)
	if s.currIndex >= lens {
		s.currIndex = 0
	}
	picked = peers[s.currIndex]
	s.currIndex++

	return
}
