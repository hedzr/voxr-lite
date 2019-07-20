module github.com/hedzr/voxr-lite

go 1.12

// replace github.com/hedzr/cmdr v0.0.0 => ../cmdr

// replace github.com/hedzr/cmdr v0.2.25 => ../cmdr

// replace github.com/hedzr/logex v0.0.0 => ../logex

replace github.com/hedzr/voxr-common v0.0.0 => ../voxr-common

replace github.com/hedzr/voxr-api v0.0.0 => ../voxr-api

require (
	github.com/Masterminds/semver v1.4.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/core v0.6.2
	github.com/go-xorm/xorm v0.7.3
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/consul/api v1.1.0
	github.com/hedzr/cmdr v1.0.1
	github.com/hedzr/logex v1.0.0
	github.com/hedzr/voxr-api v0.0.0
	github.com/hedzr/voxr-common v0.0.0
	github.com/jinzhu/gorm v1.9.9
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	go.etcd.io/etcd v3.3.13+incompatible
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
	google.golang.org/grpc v1.21.1
	gopkg.in/yaml.v2 v2.2.2
)

exclude (
	github.com/coreos/etcd v3.3.10+incompatible // indirect
	github.com/hashicorp/go-rootcerts v1.0.0
)
