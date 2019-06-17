/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"net/url"
)

type (
	Fwdr struct {
		Balancer string
		Prefix   string
		Targets  []Target `yaml:"targets,omitempty"`
	}
	Target struct {
		Url     string            `yaml:"url,omitempty"`
		Name    string            `yaml:"name,omitempty"`
		Meta    map[string]string `yaml:"meta,omitempty"`
		Rewrite map[string]string `yaml:"rewrite,omitempty"`
	}
	FwdrList struct {
		List    []Fwdr
		Rewrite map[string]string `yaml:"rewrite,omitempty"`
	}
)

var (
	fwdrList     FwdrList
	rewriteRules map[string]string
)

func check_panic(err error) {
	if err != nil {
		logrus.Fatalf("error occurs: %v", err)
	}
}

func (s *myService) OnInitForwarders(e *echo.Echo) {
	// forwarder.Init(e)

	rewriteRules = make(map[string]string)

	fObj := vxconf.GetR("server.forwarders")
	b, err := yaml.Marshal(fObj)
	check_panic(err)
	// logrus.Infof("fObj = %v", string(b))
	s2 := string(b)
	fwdrList.List = []Fwdr{}
	err = yaml.Unmarshal([]byte(s2), &fwdrList)
	check_panic(err)
	// logrus.Infof("fwdrList = %v", fwdrList)
	for i, k := range fwdrList.List {
		for j, t := range k.Targets {
			logrus.Infof("%d, %d : %v", i, j, t)
		}
	}

	for _, k := range fwdrList.List {
		if len(k.Targets) == 0 {
			continue
		}

		var bal middleware.ProxyBalancer

		var targets []*middleware.ProxyTarget
		for _, t := range k.Targets {
			u, _ := url.Parse(t.Url)
			var meta = make(echo.Map)
			for k, v := range t.Meta {
				meta[k] = v
			}
			targets = append(targets, &middleware.ProxyTarget{t.Name, u, meta})
		}
		switch k.Balancer {
		case "round-robin":
			bal = middleware.NewRoundRobinBalancer(targets)
		case "random":
			bal = middleware.NewRandomBalancer(targets)
		}

		g := e.Group(k.Prefix)
		g.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
			Balancer: bal,
			Rewrite:  fwdrList.Rewrite,
		}))
		// if len(targets[0].URL.Path) > 0 && targets[0].URL.Path != k.Prefix {
		// 	// e.Pre(middleware.RewriteWithConfig(middleware.RewriteConfig{}))
		// 	match := strings.Replace(fmt.Sprintf("%s/*", k.Prefix), "//", "/", -1)
		// 	to := strings.Replace(fmt.Sprintf("%s/$1", targets[0].URL.Path), "//", "/", -1)
		// 	rewriteRules[match] = to
		// 	logrus.Debugf("rewriting rule: %s -> %s", match, to)
		// }
	}

	// e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
}

func (s *myService) OnShutdownForwarders() {
	// forwarder.Shutdown()
}
