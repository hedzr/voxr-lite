module github.com/hedzr/voxr-lite

go 1.13

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr

// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

// replace github.com/hedzr/logex v0.0.0 => ../logex

replace github.com/hedzr/voxr-common => ../voxr-common

replace github.com/hedzr/voxr-api => ../voxr-api

// ignore github.com/Masterminds/semver

require (
	github.com/Masterminds/semver v1.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/consul/api v1.3.0
	github.com/hedzr/cmdr v1.6.18
	github.com/hedzr/errors v1.1.18
	github.com/hedzr/logex v1.1.5
	github.com/hedzr/voxr-api v0.0.0-00010101000000-000000000000
	github.com/hedzr/voxr-common v0.0.0-00010101000000-000000000000
	github.com/jinzhu/gorm v1.9.12
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/skip2/go-qrcode v0.0.0-20191027152451-9434209cb086
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	go.etcd.io/etcd v3.3.18+incompatible
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.2.7
	xorm.io/core v0.7.2
)
