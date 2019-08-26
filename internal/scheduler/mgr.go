/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

var (
	grpcHub = &GrpcHub{
		clients:         make(map[*GrpcClient]bool),
		byIds:           make(map[string]*GrpcClient),
		brand:           "inx.im",
		exited:          true,
		exiting:         make(chan bool),
		newClientAdding: make(chan *GrpcClient),
		deregister:      make(chan *GrpcClient),
		refreshClients:  make(chan string),
		invoking:        make(chan *Input),
	}

	synonym = make(map[string]string)

	// isPreferredToRealService true: 如果待替换的服务有活动实例，优先使用该活动实例，而不是使用同义词服务进行替换
	//   false: 一律替换为同义词服务
	isPreferredToRealService bool
)

// maps: key is toService, value is its replaced (synonymService)
func AddSynonyms(priorToRealService bool, maps map[string]string) {
	isPreferredToRealService = priorToRealService
	for k, v := range maps {
		synonym[k] = v
	}
}

func AddSynonym(toService, synonymService string) {
	synonym[toService] = synonymService
}

// Start starts the gRPC clients manager
func Start() echo.HandlerFunc {
	grpcHub.start()

	isPreferredToRealService = vxconf.GetBoolR("server.deps.settings.preferredToRealService", isPreferredToRealService)

	var out = make(map[string]*DepRecord)
	_ = vxconf.LoadSectionTo("server.deps", out)
	for k, v := range out {
		if !v.Disabled && !strings.EqualFold(k, "settings") {
			grpcHub.newClient(k, v)
		}
	}

	return grpcInvokeHandler
}

func Stop() {
	grpcHub.stop()
}

func RequestRefreshClient(name string) {
	grpcHub.refreshClients <- name
}

func RequestRefreshAllClients() {
	grpcHub.refreshClients <- ""
}


// newClient named a client as the value of `serviceName`.
func (h *GrpcHub) newClient(serviceName string, dr *DepRecord) (client *GrpcClient) {
	if v, ok := h.byIds[serviceName]; ok {
		return v
	}

	client = newClientInternal(serviceName, dr)

	// notify grpcHub routine the new client connection object is incoming
	h.newClientAdding <- client

	// refresh each clients after a random delay.
	go func() {
		time.Sleep(time.Duration(750+rand.Intn(750)) * time.Millisecond)
		h.refreshClients <- serviceName
	}()

	return
}

func (h *GrpcHub) start() {
	if !h.exited {
		return
	}

	go h.run()
}

func (h *GrpcHub) run() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logrus.Infof("      gRPC hub exiting.")
		ticker.Stop()
		h.reset_()
	}()

	logrus.Println("      gRPC hub started.")

	for {
		h.exited = false
		select {
		case exit := <-h.exiting:
			if exit {
				return
			}

		case tm := <-ticker.C:
			logrus.Debugf("      refreshing clients at %v", tm)
			refreshServicesSync()

		case name := <-h.refreshClients:
			if len(name) > 0 {
				logrus.Debugf("      refreshing client '%s' by signal", name)
				refreshServiceSync(name)
			}

		case input := <-h.invoking:
			// invoke_(input, grpc.EmptyCallOption{})
			logrus.Debugf("      gRPC invoking: %v", input)
			invoke__(input)

		case c := <-h.newClientAdding:
			h.addClient_(c)
			logrus.Println("=== new gRPC client session created.", h.clients) // and log it for debugging

		case c := <-h.deregister:
			h.removeClient_(c)

			// case name := <-grpcHub.querying:
			// 	//logrus.Println(name)
			// 	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*len(grpcHub.clients))*time.Second)
			// 	//defer cancel()
			// 	ctx := context.Background()
			// 	for k, ok := range grpcHub.clients {
			// 		if ok && k.conn != nil {
			// 			r, err := k.applyServiceClient.FetchApplyByUid(ctx, &pb.ApplyRequest{Uid: string(name)})
			// 			if err != nil {
			// 				logrus.Warnf("could not greet: %v", err)
			// 			} else {
			// 				logrus.Printf("Greeting: %s", r.Result)
			// 				if k.callback != nil {
			// 					k.callback(r.Result)
			// 				}
			// 			}
			// 		}
			// 	}

		}
	}
}

func (h *GrpcHub) stop() {
	if !h.exited {
		h.exiting <- true
	}
}

func (h *GrpcHub) addClient_(c *GrpcClient) {
	h.rwLock.Lock()
	defer h.rwLock.Unlock()
	h.byIds[c.Name] = c
	h.clients[c] = true // flag the client conn object is true inside map `clients`
}

func (h *GrpcHub) removeClient_(c *GrpcClient) {
	{
		h.rwLock.Lock()
		defer h.rwLock.Unlock()
		delete(h.clients, c)
		delete(h.byIds, c.Name)
	}
	h._close(c)
}

func (h *GrpcHub) _close(c *GrpcClient) {
	h.exiting <- true
	for _, p := range c.Peers {
		if p.Conn != nil {
			_ = p.Conn.Close()
			p.Conn = nil
		}
	}
}

func (h *GrpcHub) reset_() {
	h.rwLock.RLock()
	defer h.rwLock.RUnlock()
	for c, ok := range h.clients {
		if ok {
			h._close(c)
		}
	}
	h.exited = true
}

func (h *GrpcHub) broadcast_(callback func(client *GrpcClient)) {
	h.rwLock.RLock()
	defer h.rwLock.RUnlock()
	for k, ok := range h.clients {
		if ok && len(k.Peers) > 0 {
			callback(k)
		}
	}
}
