/*
 * Copyright © 2019 Hedzr Yeh.
 */

package classical

import (
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
	"net"
	"strconv"
)

func A() (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(
		"",
		grpc.WithInsecure(),
		// 负载均衡，使用 consul 作服务发现
		// grpc.WithBalancerName("round_robin"),
		grpc.WithBalancer(grpc.RoundRobin(NewConsulResolver(
			"127.0.0.1:8500", "grpc.health.v1.add",
		))),
		// https://blog.csdn.net/hatlonely/article/details/80788686
	)
	return
}

func NewConsulResolver(address string, service string) naming.Resolver {
	return &consulResolver{
		address: address,
		service: service,
	}
}

type consulResolver struct {
	address string
	service string
}

func (r *consulResolver) Resolve(target string) (naming.Watcher, error) {
	config := api.DefaultConfig()
	config.Address = r.address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &consulWatcher{
		client:  client,
		service: r.service,
		addrs:   map[string]struct{}{},
	}, nil
}

type consulWatcher struct {
	client    *api.Client
	service   string
	addrs     map[string]struct{}
	lastIndex uint64
}

func (w *consulWatcher) Next() ([]*naming.Update, error) {
	for {
		services, metainfo, err := w.client.Health().Service(w.service, "", true, &api.QueryOptions{
			WaitIndex: w.lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
		})
		if err != nil {
			logrus.Warn("error retrieving instances from Consul: %v", err)
		}
		w.lastIndex = metainfo.LastIndex

		addrs := map[string]struct{}{}
		for _, service := range services {
			addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
		}

		var updates []*naming.Update
		for addr := range w.addrs {
			if _, ok := addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr})
			}
		}

		for addr := range addrs {
			if _, ok := w.addrs[addr]; !ok {
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
			}
		}

		if len(updates) != 0 {
			w.addrs = addrs
			return updates, nil
		}
	}
}

func (w *consulWatcher) Close() {
}
