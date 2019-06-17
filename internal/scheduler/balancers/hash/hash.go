/*
 * Copyright © 2019 Hedzr Yeh.
 */

package hash

import (
	"errors"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"hash/crc32"
	"math/rand"
	"strconv"
)

type (
	HashBalancer struct {
		key string
	}
)

func New() balance2.Balancer {
	return &HashBalancer{""}
}

func (s *HashBalancer) Name() string {
	return balance2.HASH
}

func (s *HashBalancer) Clone() balance2.Balancer {
	return &HashBalancer{s.key}
}

func (s *HashBalancer) UpdateScales(scales map[string]int) (b balance2.Balancer) {
	b = s
	return
}

func (s *HashBalancer) Pick(peers []*balance2.ServicePeer) (picked *balance2.ServicePeer, err error) {
	if len(peers) == 0 {
		err = errors.New("no instance")
		return
	}

	defKey := strconv.Itoa(rand.Int())

	// if len(key) > 0 {
	// 	defKey = key[0]
	// }

	lens := len(peers)
	hashVal := crc32.Checksum([]byte(defKey), crc32.MakeTable(crc32.IEEE))
	index := int(hashVal) % lens
	picked = peers[index]

	return
}
