module github.com/makerdao/vulcanizedb

go 1.12

require (
	github.com/dave/jennifer v1.3.0
	github.com/ethereum/go-ethereum v1.9.8
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hpcloud/tail v1.0.0
	github.com/jmoiron/sqlx v0.0.0-20181024163419-82935fac6c1a
	github.com/lib/pq v1.0.0
	github.com/makerdao/vdb-mcd-transformers v0.2.17
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.10.0
	github.com/pressly/goose v2.6.0+incompatible
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
)

replace github.com/ethereum/go-ethereum => github.com/makerdao/go-ethereum v1.9.15-statechange-filter

replace gopkg.in/urfave/cli.v1 => gopkg.in/urfave/cli.v1 v1.20.0
