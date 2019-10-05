module github.com/hedzr/voxr-lite

go 1.12

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr

// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

// replace github.com/hedzr/logex v0.0.0 => ../logex

replace github.com/hedzr/voxr-common v0.0.0 => ../voxr-common

replace github.com/hedzr/voxr-api v0.0.0 => ../voxr-api

require (
	github.com/Masterminds/semver v1.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/consul/api v1.2.0
	github.com/hedzr/cmdr v1.5.3
	github.com/hedzr/logex v1.0.3
	github.com/hedzr/voxr-api v0.0.0
	github.com/hedzr/voxr-common v0.0.0
	github.com/jinzhu/gorm v1.9.11
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	go.etcd.io/etcd v3.3.15+incompatible
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/net v0.0.0-20191003171128-d98b1b443823
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
	xorm.io/core v0.7.2-0.20190928055935-90aeac8d08eb
)
