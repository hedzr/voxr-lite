/*
 * Copyright © 2019 Hedzr Yeh.
 */

package scheduler

import (
	"github.com/golang/protobuf/proto"
	"sync"
	"time"
)

const (
	pingPeriod = 2 * 60 * time.Second
)

type (
	DepRecord struct {
		Id       string `yaml:"id"`
		Addr     string `yaml:"addr"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Disabled bool   `yaml:"disabled"`
		Balancer BC     `yaml:"balancer"`
	}

	BC struct {
		Type     string         `yaml:"type"`
		SubType  string         `yaml:"sub-type"`
		Versions map[string]int `yaml:"versions"`
	}

	Input struct {
		ServiceName   string
		PBServiceName string
		MethodName    string
		In            proto.Message
		Out           proto.Message
		Callback      func(error, *Input, proto.Message)
		client        *GrpcClient
	}

	GrpcHub struct {
		clients         map[*GrpcClient]bool
		byIds           map[string]*GrpcClient
		rwLock          sync.RWMutex
		brand           string
		exited          bool
		exiting         chan bool
		newClientAdding chan *GrpcClient
		deregister      chan *GrpcClient
		refreshClients  chan string
		invoking        chan *Input
	}
)
