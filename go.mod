module github.com/vulcanize/vulcanizedb

go 1.12

require (
	bazil.org/fuse v0.0.0-20180421153158-65cc252bf669
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9
	github.com/Stebalien/go-bitfield v0.0.1
	github.com/allegro/bigcache v0.0.0-20190618191010-69ea0af04088
	github.com/aristanetworks/goarista v0.0.0-20190712234253-ed1100a1c015
	github.com/bren2010/proquint v0.0.0-20160323162903-38337c27106d
	github.com/btcsuite/btcd v0.0.0-20190629003639-c26ffa870fd8
	github.com/btcsuite/goleveldb v1.0.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v3 v3.0.0
	github.com/cheekybits/genny v1.0.0
	github.com/coreos/go-semver v0.3.0
	github.com/cskr/pubsub v1.0.2
	github.com/dave/jennifer v1.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/davidlazar/go-crypto v0.0.0-20190522120613-62389b5e4ae0
	github.com/deckarep/golang-set v1.7.1
	github.com/dgraph-io/badger/v2 v2.0.0-rc.2+incompatible
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2
	github.com/dustin/go-humanize v1.0.0
	github.com/edsrzf/mmap-go v1.0.0
	github.com/elastic/gosigar v0.10.4
	github.com/ethereum/go-ethereum v1.9.1
	github.com/facebookgo/atomicfile v0.0.0-20151019160806-2de1f203e7d5
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-stack/stack v1.8.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/golang-lru v0.5.3
	github.com/hashicorp/hcl v1.0.0
	github.com/hpcloud/tail v1.0.0
	github.com/huin/goupnp v1.0.0
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/ipfs/bbloom v0.0.1
	github.com/ipfs/go-bitswap v0.1.6
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-blockservice v0.1.2
	github.com/ipfs/go-cid v0.0.3
	github.com/ipfs/go-cidutil v0.0.2
	github.com/ipfs/go-datastore v0.0.5
	github.com/ipfs/go-ds-badger v0.0.5
	github.com/ipfs/go-ds-flatfs v0.0.2
	github.com/ipfs/go-ds-leveldb v0.0.2
	github.com/ipfs/go-ds-measure v0.0.1
	github.com/ipfs/go-fs-lock v0.0.1
	github.com/ipfs/go-ipfs v0.4.22
	github.com/ipfs/go-ipfs-blockstore v0.0.1
	github.com/ipfs/go-ipfs-chunker v0.0.1
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/go-ipfs-delay v0.0.1
	github.com/ipfs/go-ipfs-ds-help v0.0.1
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipfs-exchange-offline v0.0.1
	github.com/ipfs/go-ipfs-files v0.0.3
	github.com/ipfs/go-ipfs-posinfo v0.0.1
	github.com/ipfs/go-ipfs-pq v0.0.1
	github.com/ipfs/go-ipfs-routing v0.1.0
	github.com/ipfs/go-ipfs-util v0.0.1
	github.com/ipfs/go-ipld-cbor v0.0.3
	github.com/ipfs/go-ipld-format v0.0.2
	github.com/ipfs/go-ipld-git v0.0.2
	github.com/ipfs/go-ipns v0.0.1
	github.com/ipfs/go-log v0.0.1
	github.com/ipfs/go-merkledag v0.1.0
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-mfs v0.1.1
	github.com/ipfs/go-path v0.0.7
	github.com/ipfs/go-peertaskqueue v0.1.1
	github.com/ipfs/go-todocounter v0.0.1
	github.com/ipfs/go-unixfs v0.1.0
	github.com/ipfs/go-verifcid v0.0.1
	github.com/ipfs/interface-go-ipfs-core v0.1.0
	github.com/jackpal/gateway v1.0.5
	github.com/jackpal/go-nat-pmp v1.0.1
	github.com/jbenet/go-is-domain v1.0.2
	github.com/jbenet/go-temp-err-catcher v0.0.0-20150120210811-aac704a3f4f2
	github.com/jbenet/goprocess v0.1.3
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jmoiron/sqlx v0.0.0-20190426154859-38398a30ed85
	github.com/karalabe/usb v0.0.0-20190703133951-9be757f914c0
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/kkdai/bstream v0.0.0-20181106074824-b3251f7901ec // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/koron/go-ssdp v0.0.0-20180514024734-4a0ed625a78b
	github.com/lib/pq v1.2.0
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-buffer-pool v0.0.2
	github.com/libp2p/go-conn-security-multistream v0.1.0
	github.com/libp2p/go-eventbus v0.0.3
	github.com/libp2p/go-flow-metrics v0.0.1
	github.com/libp2p/go-libp2p v0.1.2
	github.com/libp2p/go-libp2p-autonat v0.1.0
	github.com/libp2p/go-libp2p-autonat-svc v0.1.0
	github.com/libp2p/go-libp2p-circuit v0.1.0
	github.com/libp2p/go-libp2p-connmgr v0.1.0
	github.com/libp2p/go-libp2p-core v0.0.9
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-discovery v0.1.0
	github.com/libp2p/go-libp2p-kad-dht v0.1.1
	github.com/libp2p/go-libp2p-kbucket v0.2.0
	github.com/libp2p/go-libp2p-loggables v0.1.0
	github.com/libp2p/go-libp2p-mplex v0.2.1
	github.com/libp2p/go-libp2p-nat v0.0.4
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/libp2p/go-libp2p-pnet v0.1.0
	github.com/libp2p/go-libp2p-protocol v0.1.0 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.1.0
	github.com/libp2p/go-libp2p-pubsub-router v0.1.0
	github.com/libp2p/go-libp2p-quic-transport v0.1.1
	github.com/libp2p/go-libp2p-record v0.1.0
	github.com/libp2p/go-libp2p-routing v0.1.0
	github.com/libp2p/go-libp2p-routing-helpers v0.1.0
	github.com/libp2p/go-libp2p-secio v0.1.1
	github.com/libp2p/go-libp2p-swarm v0.1.1
	github.com/libp2p/go-libp2p-tls v0.1.0
	github.com/libp2p/go-libp2p-transport-upgrader v0.1.1
	github.com/libp2p/go-libp2p-yamux v0.2.1
	github.com/libp2p/go-maddr-filter v0.0.5
	github.com/libp2p/go-mplex v0.1.0
	github.com/libp2p/go-msgio v0.0.4
	github.com/libp2p/go-nat v0.0.3
	github.com/libp2p/go-reuseport v0.0.1
	github.com/libp2p/go-reuseport-transport v0.0.2
	github.com/libp2p/go-stream-muxer v0.1.0 // indirect
	github.com/libp2p/go-stream-muxer-multistream v0.2.0
	github.com/libp2p/go-tcp-transport v0.1.0
	github.com/libp2p/go-testutil v0.1.0 // indirect
	github.com/libp2p/go-ws-transport v0.1.0
	github.com/libp2p/go-yamux v1.2.3
	github.com/lucas-clemente/quic-go v0.11.2
	github.com/magiconair/properties v1.8.1
	github.com/marten-seemann/qtls v0.2.4
	github.com/mattn/go-colorable v0.1.2
	github.com/mattn/go-isatty v0.0.8
	github.com/mattn/go-runewidth v0.0.4
	github.com/miekg/dns v1.1.15
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/minio/sha256-simd v0.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mmcloughlin/avo v0.0.0-20190731014047-bb615f61ce85
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-base32 v0.0.3
	github.com/multiformats/go-multiaddr v0.0.4
	github.com/multiformats/go-multiaddr-dns v0.0.3
	github.com/multiformats/go-multiaddr-fmt v0.0.1
	github.com/multiformats/go-multiaddr-net v0.0.1
	github.com/multiformats/go-multibase v0.0.1
	github.com/multiformats/go-multihash v0.0.6
	github.com/multiformats/go-multistream v0.1.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.4.0
	github.com/pkg/errors v0.8.1
	github.com/polydawn/refmt v0.0.0-20190731040541-eff0b363297a
	github.com/pressly/goose v2.6.0+incompatible
	github.com/prometheus/tsdb v0.10.0
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spacemonkeygo/openssl v0.0.0-20181017203307-c2dcc5cca94a
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/steakknife/bloomfilter v0.0.0-20180906043351-99ee86d9200f
	github.com/steakknife/hamming v0.0.0-20180906055317-003c143a81c2
	github.com/syndtr/goleveldb v1.0.0
	github.com/tyler-smith/go-bip39 v1.0.0
	github.com/vulcanize/eth-block-extractor v0.0.0-20190801172153-2835f21156aa
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc
	github.com/whyrusleeping/chunker v0.0.0-20181014151217-fe64bd25879f
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1
	github.com/whyrusleeping/go-logging v0.0.0-20170515211332-0457bb6b88fc
	github.com/whyrusleeping/go-notifier v0.0.0-20170827234753-097c5d47330f
	github.com/whyrusleeping/mafmt v1.2.8
	github.com/whyrusleeping/mdns v0.0.0-20180901202407-ef14215e6b30
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/whyrusleeping/sql-datastore v0.0.0-20190124195324-b24eb8d0ce14
	github.com/whyrusleeping/timecache v0.0.0-20160911033111-cfcb2f1abfee
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208
	go.opencensus.io v0.22.0
	go.uber.org/atomic v1.4.0
	go.uber.org/dig v1.7.0
	go.uber.org/fx v1.9.0
	go.uber.org/multierr v1.1.0
	go4.org v0.0.0-20190313082347-94abd6928b1d
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190801041406-cbf593c0f2f3
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190802003818-e9bb7d36c060
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/dgraph-io/badger v2.0.0-rc.2+incompatible => github.com/dgraph-io/badger/v2 v2.0.0-rc.2

replace github.com/dgraph-io/badger/v2 v2.0.0-rc2 => github.com/dgraph-io/badger v1.6.0-rc1

replace github.com/ipfs/go-ipfs v0.4.22 => github.com/vulcanize/go-ipfs v0.4.22-alpha

replace github.com/ipfs/go-ipfs-config v0.0.3 => github.com/vulcanize/go-ipfs-config v0.0.8-alpha

replace github.com/ethereum/go-ethereum v1.9.1 => github.com/vulcanize/go-ethereum v0.0.0-20190731183759-8e20673bd101
