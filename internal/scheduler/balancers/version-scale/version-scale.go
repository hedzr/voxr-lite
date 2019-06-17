/*
 * Copyright © 2019 Hedzr Yeh.
 */

package version_scale

import (
	"errors"
	"github.com/Masterminds/semver"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type (
	VersionScaleBalancer struct {
		// Servers    []*balance.ServicePeer
		// lastPicked int
		subType        string
		currIndex      int
		rrIndex        int
		wheel          []int
		wheelKeys      []*semver.Constraints
		wheelRrIndexes []int
		scales         map[string]int // key: ver-matcher such as '>=1.1.3'
	}
)

func New(scales map[string]int) (b balance2.Balancer) {
	b = (&VersionScaleBalancer{balance2.RANDOM, 0, 0, nil, nil, nil, nil}).UpdateScales(scales)
	return
}

func (s *VersionScaleBalancer) Name() string {
	return balance2.VERSIONSCALE
}

func (s *VersionScaleBalancer) Clone() balance2.Balancer {
	return &VersionScaleBalancer{s.subType, s.currIndex, s.rrIndex, s.wheel, s.wheelKeys, s.wheelRrIndexes, s.scales}
}

func (s *VersionScaleBalancer) UpdateScales(scales map[string]int) (b balance2.Balancer) {
	b = s

	if len(scales) == 0 {
		return
	}

	changed := false
	if len(s.scales) != 0 {
		for k, v := range scales {
			if z, ok := s.scales[k]; ok {
				if z != v {
					changed = true
					break
				}
			}
		}
	} else {
		changed = true
	}
	if !changed {
		return
	}

	s.scales = scales
	s.wheel = nil
	s.wheelKeys = nil

	for k, v := range scales {
		c, err := semver.NewConstraint(k)
		if err != nil {
			logrus.Fatalf("Version Constraint Error: %v", err)
		}

		s.wheelRrIndexes = append(s.wheelRrIndexes, 0)
		s.wheelKeys = append(s.wheelKeys, c)
		i := len(s.wheelKeys) - 1
		for x := 0; x < v; x++ {
			s.wheel = append(s.wheel, i)
		}
	}

	dest := make([]int, len(s.wheel))
	perm := rand.Perm(len(s.wheel))
	for i, v := range perm {
		dest[v] = s.wheel[i]
	}
	s.wheel = dest

	return
}

func (s *VersionScaleBalancer) Pick(peers []*balance2.ServicePeer) (picked *balance2.ServicePeer, err error) {
	if len(peers) == 0 {
		err = errors.New("no instance")
		return
	}

	// sub-type == round-robin on a shuffled wheel

	if s.currIndex >= len(s.wheel) {
		s.currIndex = 0
	}
	pickedVmI := s.wheel[s.currIndex]
	pickedVm := s.wheelKeys[pickedVmI]
	s.currIndex++

	var candicated []*balance2.ServicePeer
	for _, p := range peers {
		v, err := semver.NewVersion(p.Record.Version)
		if err != nil {
			logrus.Warnf("invalid version number: %v, error: %v", p.Record.Version, err)
		}
		if p != nil && p.Record != nil && pickedVm.Check(v) {
			candicated = append(candicated, p)
		}
	}

	if len(candicated) == 0 {
		s.rrIndex++
		s.rrIndex %= len(peers)
		picked = peers[s.rrIndex]
	} else {
		s.wheelRrIndexes[pickedVmI]++
		s.wheelRrIndexes[pickedVmI] %= len(candicated)
		picked = candicated[s.wheelRrIndexes[pickedVmI]]
	}
	return
}
