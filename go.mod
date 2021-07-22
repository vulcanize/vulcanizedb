module github.com/vulcanize/vulcanizedb

go 1.12

require (
	github.com/aristanetworks/goarista v0.0.0-20190712234253-ed1100a1c015 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/dave/jennifer v1.3.0
	github.com/davidlazar/go-crypto v0.0.0-20190522120613-62389b5e4ae0 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/elastic/gosigar v0.10.4 // indirect
	github.com/ethereum/go-ethereum v1.9.1
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hpcloud/tail v1.0.0
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-blockservice v0.1.4
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-ipfs v0.8.0
	github.com/ipfs/go-ipfs-blockstore v0.1.4
	github.com/ipfs/go-ipfs-ds-help v0.1.1
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/jmoiron/sqlx v0.0.0-20190426154859-38398a30ed85
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.2.0
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multihash v0.0.14
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pressly/goose v2.6.0+incompatible
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
)

replace github.com/dgraph-io/badger v2.0.0-rc.2+incompatible => github.com/dgraph-io/badger/v2 v2.0.0-rc.2

replace github.com/dgraph-io/badger/v2 v2.0.0-rc2 => github.com/dgraph-io/badger v1.6.0-rc1

replace github.com/ipfs/go-ipfs v0.4.22 => github.com/vulcanize/go-ipfs v0.4.22-alpha

replace github.com/ipfs/go-ipfs-config v0.0.3 => github.com/vulcanize/go-ipfs-config v0.0.8-alpha

replace github.com/ethereum/go-ethereum v1.9.1 => github.com/vulcanize/go-ethereum v1.9.11-statediff-0.0.2
