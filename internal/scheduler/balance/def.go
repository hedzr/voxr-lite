/*
 * Copyright © 2019 Hedzr Yeh.
 */

package balance

import (
	"fmt"
	"github.com/hedzr/voxr-common/kvs/store"
	"google.golang.org/grpc"
)

const (
	RANDOM       = "random"
	ROUNDROBIN   = "round-robin"
	VERSIONSCALE = "weighted-version"
	HASH         = "hash"
)

type (
	// Never Used.
	SvcRec interface {
	}

	ServicePeer struct {
		Record *store.ServiceRecord
		Conn   *grpc.ClientConn
		Meta   map[string]interface{}
	}

	Balancer interface {
		Name() string
		Pick(peers []*ServicePeer) (picked *ServicePeer, err error)
		UpdateScales(scale map[string]int) Balancer
		Clone() Balancer
	}

	Balancers map[string]Balancer
)

func (sp *ServicePeer) String() (s string) {
	return fmt.Sprintf("%v:%v (%s) v%v", sp.Record.IP, sp.Record.Port, sp.Record.ID, sp.Record.Version)
}

func EmptyPeer() *ServicePeer {
	return &ServicePeer{}
}

func NewPeer(r *store.ServiceRecord, c *grpc.ClientConn) *ServicePeer {
	return &ServicePeer{r, c, make(map[string]interface{})}
}
