/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"fmt"
	"github.com/hedzr/voxr-api/api"
	"github.com/hedzr/voxr-common/kvs/store"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/vxconf/gwk"
	"github.com/hedzr/voxr-lite/internal/scheduler/balance"
	"github.com/hedzr/voxr-lite/internal/scheduler/balancers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

func refreshServicesSync() {
	grpcHub.rwLock.RLock()
	defer grpcHub.rwLock.RUnlock()
	for _, c := range grpcHub.byIds {
		refreshClient(c)
	}

	debugDump()
}

func refreshServiceSync(serviceName string) {
	if len(serviceName) == 0 {
		refreshServicesSync()
		return
	}

	grpcHub.rwLock.RLock()
	defer grpcHub.rwLock.RUnlock()
	if c, ok := grpcHub.byIds[serviceName]; ok {
		// grpcHub.rwLock.RUnlock()
		refreshClient(c)
		debugDump()
	}
}

func refreshClient(c *GrpcClient) {
	records := gwk.ThisConfig.Registrar.SvrRecordResolverAll(c.Name, api.GRPC)
	// if c.Name == "vx-core" {
	// 	logrus.Debugf("debug: vx-core: %v", records)
	// }
	if len(records) == 0 {
		var px *balance.ServicePeer
		for _, p := range c.Peers {
			if p.Record.IsLocalDefined() {
				px = p
				break
			}
		}
		c.Peers = make([]*balance.ServicePeer, 0)
		if px != nil {
			c.Peers = append(c.Peers, px)
		}
		return
	}

	var todo []int
	var todoRecords []*balance.ServicePeer
	for z, p := range c.Peers {
		found := false
		for i, r := range records {
			if p.Record.Equal(r) {
				records = append(records[:i], records[i+1:]...)
				found = true
				break
			}
		}
		if !found && !p.Record.IsLocalDefined() {
			todo = append(todo, z)
			todoRecords = append(todoRecords, p)
			logrus.Debugf("    - [%s] adding the invalid peer: %v", c.Name, p.Record)
		}
	}

	for i := len(todo) - 1; i >= 0; i-- {
		z := todo[i]
		c.Peers = append(c.Peers[0:z], c.Peers[z+1:]...)
		p := todoRecords[i]
		if p.Conn != nil {
			p.Conn.Close()
		}
		logrus.Debugf("    - [%s] remove the invalid peer: %v", c.Name, p.Record)
	}

	for _, r := range records {
		addr := fmt.Sprintf("%v:%d", r.IP, r.Port)
		conn, err := grpc.Dial(addr,
			grpc.WithInsecure(),
			// grpc.WithBalancerName(grpcBalancerName),
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(keepAlive))
		if err != nil {
			logrus.Fatalf("did not connect: %v", err)
		} else {
			c.Peers = append(c.Peers, balance.NewPeer(r, conn))
			// c.balancer.UpdateScales()

			logrus.Debugf("        - [%s] add new peer: %v", c.Name, r)
		}
	}
}

const grpcBalancerName = "round-robin"

var keepAlive = keepalive.ClientParameters{
	10 * time.Second,
	20 * time.Second,
	true,
}

func lookupClients(serviceName string, dr *DepRecord) (client *GrpcClient) {
	addrLocal := tool.LoadHostDefinition(fmt.Sprintf("server.deps.%s", serviceName))
	addrLocal = dr.Addr
	if len(dr.Host) > 0 && dr.Port > 0 {
		addrLocal = fmt.Sprintf("%v:%v", dr.Host, dr.Port)
	}

	// balancers.RegisterResolver(fmt.Sprintf("voxr/services/%s/peers", serviceName))

	client = &GrpcClient{
		serviceName,
		make([]*balance.ServicePeer, 0),
		make(chan *Input, 8),
		make(chan bool),
		true,
		dr,
		balancers.NewWith(dr.Balancer.Type, dr.Balancer.SubType, dr.Balancer.Versions),
	}

	// X
	records := gwk.ThisConfig.Registrar.SvrRecordResolverAll(serviceName, api.GRPC)
	if len(records) > 0 {
		// ix := rand.Intn(len(records))
		// addr = fmt.Sprintf("%v:%d", records[ix].IP, records[ix].Port)
		for _, r := range records {
			addr := fmt.Sprintf("%v:%d", r.IP, r.Port)
			if addr == addrLocal {
				addrLocal = ""
			}
			logrus.Debugf("    - [%s] connecting to: %v", serviceName, addr)
			// Set up a connection to the server.
			conn, err := grpc.Dial(addr,
				grpc.WithInsecure(),
				// grpc.WithBalancerName(grpcBalancerName),
				// grpc.WithBlock(),
				grpc.WithKeepaliveParams(keepAlive))
			if err != nil {
				logrus.Fatalf("did not connect: %v", err)
			} else {
				// defer conn.Close()
				client.Peers = append(client.Peers, balance.NewPeer(r, conn))
				logrus.Debugf("    - [%s] add new peer: %v", serviceName, r)
			}
		}
	}

	// local static record
	if len(addrLocal) > 0 {
		conn, err := grpc.Dial(addrLocal, grpc.WithInsecure())
		if err != nil {
			logrus.Fatalf("did not connect: %v", err)
		}

		client.Peers = append(client.Peers, balance.NewPeer(store.NewServiceRecord(addrLocal), conn))
		logrus.Debugf("    - [%s] add new peer: %v", serviceName, addrLocal)
	}

	debugDump()
	return
}

func debugDump() {
	if vxconf.GetBoolR("server.deps-debug", false) {
		logrus.Debug("    -------- debug of grpc services:")
		for _, c := range grpcHub.byIds {
			for _, p := range c.Peers {
				logrus.Debugf("    - [%s] peer: %v", c.Name, p.Record)
			}
		}
	}
}
