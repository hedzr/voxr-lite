/*
 * Copyright © 2019 Hedzr Yeh.
 */

package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServerConfig struct {
	// NEVER USED
	Port int
}

type Config struct {
	// NEVER USED
	Server ServerConfig
	Db     DatabaseConfig
}

var ServerPort int // NEVER USED

func InitConfig() *Config { // NEVER USED
	data, _ := ioutil.ReadFile("config.yml")
	fmt.Println(string(data))
	config := Config{}
	yaml.UnmarshalStrict(data, &config)
	ServerPort = config.Server.Port
	return &config
}
