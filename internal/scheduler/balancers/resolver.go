/*
 * Copyright © 2019 Hedzr Yeh.
 */

package balancers

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"strings"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "2379"
)

var (
	defaultMinFrequency = 120 * time.Second
)

func init() {
}

type etcdBuilder struct {
	watchKeyPrefix string
}

func NewETCDBuilder() resolver.Builder {
	return &etcdBuilder{}
}

func RegisterResolver(keyPrefix string) {
	resolver.Register(&etcdBuilder{watchKeyPrefix: keyPrefix})
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	etcdProxys, err := parseTarget(target.Endpoint)
	if err != nil {
		return nil, err
	}

	grpclog.Infoln("etcd resolver, endpoints:", etcdProxys)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdProxys,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, errors.New("connect to etcd proxy error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	rlv := &etcdResolver{
		cc:             cc,
		cli:            cli,
		ctx:            ctx,
		cancel:         cancel,
		watchKeyPrefix: b.watchKeyPrefix,
		freq:           5 * time.Second,
		t:              time.NewTimer(0),
		rn:             make(chan struct{}, 1),
		im:             make(chan []resolver.Address),
		wg:             sync.WaitGroup{},
	}

	rlv.wg.Add(2)
	go rlv.watcher()
	go rlv.FetchBackendsWithWatch()

	return rlv, nil
}

type etcdResolver struct {
	retry  int
	freq   time.Duration
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	cli    *clientv3.Client
	t      *time.Timer

	watchKeyPrefix string

	rn chan struct{}
	im chan []resolver.Address

	wg sync.WaitGroup
}

func (r *etcdResolver) ResolveNow(opt resolver.ResolveNowOption) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

func (r *etcdResolver) Close() {
	r.cancel()
	r.wg.Wait()
	r.t.Stop()
}

func (r *etcdResolver) watcher() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case addrs := <-r.im:
			if len(addrs) > 0 {
				r.retry = 0
				r.t.Reset(r.freq)
				r.cc.NewAddress(addrs)
				continue
			}
		case <-r.t.C:
		case <-r.rn:
		}

		result := r.FetchBackends()

		if len(result) == 0 {
			r.retry++
			r.t.Reset(r.freq)
		} else {
			r.retry = 0
			r.t.Reset(r.freq)
		}

		r.cc.NewAddress(result)
	}
}

func (r *etcdResolver) FetchBackendsWithWatch() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case _ = <-r.cli.Watch(r.ctx, r.watchKeyPrefix, clientv3.WithPrefix()):
			result := r.FetchBackends()
			r.im <- result
		}
	}
}

func (r *etcdResolver) FetchBackends() []resolver.Address {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := make([]resolver.Address, 0)

	resp, err := r.cli.Get(ctx, r.watchKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		grpclog.Errorln("Fetch etcd proxy error:", err)
		return result
	}

	for _, kv := range resp.Kvs {
		if strings.TrimSpace(string(kv.Value)) == "" {
			continue
		}
		result = append(result, resolver.Address{Addr: string(kv.Value)})
	}

	grpclog.Infoln(">>>>> endpoints fetch: ", result)

	return result
}

func parseTarget(target string) ([]string, error) {
	var (
		endpoints = make([]string, 0)
	)

	if target == "" {
		return nil, errors.New("invalid target")
	}

	for _, endpoint := range strings.Split(target, ",") {
		if ip := net.ParseIP(endpoint); ip != nil {
			endpoints = append(endpoints, net.JoinHostPort(endpoint, defaultPort))
			continue
		}

		if _, port, err := net.SplitHostPort(endpoint); err == nil {
			if port == "" {
				return endpoints, errors.New("Invalid address format")
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}
