/*
 * Copyright © 2019 Hedzr Yeh.
 */

package balancers_test

import (
	"github.com/hedzr/voxr-common/kvs/store"
	balance2 "github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/hash"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/rand"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/round-robin"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers/version-scale"
	"net"
	"testing"
)

var (
	peers = []*balance2.ServicePeer{
		{&store.ServiceRecord{
			net.ParseIP("10.9.1.121"), 3500, "AA", "1.1.0", "grpc", false,
		}, nil, make(map[string]interface{})},
		{&store.ServiceRecord{
			net.ParseIP("10.9.1.122"), 3501, "BB", "1.1.0", "grpc", false,
		}, nil, make(map[string]interface{})},
		{&store.ServiceRecord{
			net.ParseIP("10.9.1.123"), 3502, "CC", "1.2.1", "grpc", false,
		}, nil, make(map[string]interface{})},
		// {&store.ServiceRecord{
		// 	net.ParseIP("10.9.1.217"), 3502, "DD", "1.2.9", "grpc", false,
		// }, nil, make(map[string]interface{})},
	}
)

const (
	rounds = 800
)

func TestForRand(t *testing.T) {
	b := rand.New()
	tester(b, t)
}

func TestForRoundRobin(t *testing.T) {
	b := round_robin.New()
	tester(b, t)
}

func TestForHash(t *testing.T) {
	b := hash.New()
	tester(b, t)
}

func TestForVersionScale(t *testing.T) {
	b := version_scale.New(map[string]int{
		"1.1.x": 80,
		"1.2.x": 20,
	})
	tester(b, t)
}

func tester(b balance2.Balancer, t *testing.T) {
	results := make(map[*store.ServiceRecord]int)

	for i := 0; i < rounds; i++ {
		peer, err := b.Pick(peers)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		results[peer.Record]++
		// t.Logf("    %-5d. picked: %v", i, peer.String())
	}

	t.Log("*** Results:")
	for k, v := range results {
		t.Logf("    %v: %d", k, v)
	}
	t.Log("OK")
}
