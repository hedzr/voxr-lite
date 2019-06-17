/*
 * Copyright © 2019 Hedzr Yeh.
 */

package r2

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"github.com/golang/protobuf/proto"
// 	"net"
// 	"os"
// 	"strings"
// 	"time"
//
// 	"go.etcd.io/etcd/clientv3"
// 	"google.golang.org/grpc"
// )
//
// func init() {
// 	go RPCServeForever()
// }
//
// func RPCServeForever() error {
// 	var (
// 		err      error
// 		listener net.Listener
// 	)
//
// 	srv := grpc.NewServer()
//
// 	proto.RegisterQcloudServiceServer(srv, new(QCloudRpcServer))
//
// 	if listener, err = net.Listen("tcp4", util.AppConfig.RpcListen); err != nil {
// 		return err
// 	}
//
// 	fmt.Println("[D] rpc will serve on %s", util.AppConfig.RpcListen)
//
// 	/* stupid but useful */
// 	go registerMyself(5)
//
// 	return srv.Serve(listener)
// }
//
// func registerMyself(ttl int64) {
// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   util.AppConfig.ETCDCluster,
// 		DialTimeout: 3 * time.Second,
// 	})
//
// 	if err != nil {
// 		util.BhAlarm(util.BH_LOG_SYSTEM, err, "New etcd client error")
// 		panic(err)
// 	}
//
// 	_myself := myself()
// 	resp := NewLeaseGrant(cli, _myself, ttl*2)
//
// 	stop := make(chan struct{})
// 	util.RegistExitHook(func() error {
// 		stop <- struct{}{}
// 		if cli == nil {
// 			return nil
// 		}
//
// 		defer cli.Close()
//
// 		if delResp, err := cli.Delete(context.TODO(), _myself); err == nil {
// 			fmt.Println("sai yo na ra:", delResp)
// 		}
//
// 		return err
// 	})
//
// 	t := time.NewTicker(time.Duration(ttl) * time.Second)
//
// 	for {
// 		select {
// 		case <-t.C:
// 			kres, e := cli.KeepAliveOnce(context.TODO(), resp.ID)
// 			if e != nil {
// 				fmt.Fprintf(os.stderr, "[E] keepalive_once error: %s", e.Error())
// 				// bugfix: 重新注册
// 				// 某次观察到etcd集群工作正常，但keepAliveOnce一直维持不住导致服务不可用
// 				resp = NewLeaseGrant(cli, _myself, 2*ttl)
// 			} else {
// 				fmt.Println("[D] keepalive response, version:", kres.Revision, "raft_term:", kres.RaftTerm)
// 			}
// 		case <-stop:
// 			t.Stop()
// 			fmt.Println("Oops: goodbye")
// 			return
// 		}
// 	}
//
// 	return
// }
//
// func myself() string {
// 	hostName, _ := os.Hostname()
// 	return `/xxx_service/` + hostName
// }
//
// // FIXME: the first one ?
// func getLocalIPV4Addr(port string) string {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	for _, addr := range addrs {
// 		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
// 			if ip.IP.To4() == nil {
// 				continue
// 			}
//
// 			if strings.HasPrefix(port, ":") {
// 				return ip.IP.String() + port
// 			}
//
// 			return net.JoinHostPort(ip.IP.String(), port)
// 		}
// 	}
//
// 	return ""
// }
//
// func NewLeaseGrant(client *clientv3.Client, value string, ttl int64) *clientv3.LeaseGrantResponse {
// 	if client == nil {
// 		panic(errors.New("Invalid etcd client"))
// 	}
//
// 	resp, err := client.Grant(context.TODO(), ttl*2) // must longer
// 	if err != nil {
// 		util.BhAlarm(util.BH_LOG_SYSTEM, err, "Grant ucenter error")
// 		panic(err)
// 	}
//
// 	_, err = client.Put(context.TODO(), value, getLocalIPV4Addr(util.AppConfig.RpcListen), clientv3.WithLease(resp.ID))
// 	if err != nil {
// 		util.BhAlarm(util.BH_LOG_SYSTEM, err, "Put myself into etcd cluster error")
// 		panic(err)
// 	}
//
// 	return resp
// }
