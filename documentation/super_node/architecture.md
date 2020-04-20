# VulcanizeDB Super Node Architecture
The VulcanizeDB super node is a collection of interfaces that are used to extract, process, and store in Postgres-IPFS
all chain data. The raw data indexed by the super node serves as the basis for more specific watchers and applications.

Currently the service supports complete processing of all Bitcoin and Ethereum data.

## Table of Contents
1. [Processes](#processes)
1. [Command](#command)
1. [Configuration](#config)
1. [Database](#database)
1. [APIs](#apis)
1. [Resync](#resync)
1. [IPFS Considerations](#ipfs-considerations)

## Processes
The [super node service](../../pkg/super_node/service.go#L61) is a watcher comprised of the following interfaces:

* [Payload Fetcher](../../pkg/super_node/shared/interfaces.go#L29): Fetches raw chain data from a half-duplex endpoint (HTTP/IPC), used for historical data fetching. ([BTC](../../pkg/super_node/btc/payload_fetcher.go), [ETH](../../pkg/super_node/eth/payload_fetcher.go)).
* [Payload Streamer](../../pkg/super_node/shared/interfaces.go#L24): Streams raw chain data from a full-duplex endpoint (WebSocket/IPC), used for syncing data at the head of the chain in real-time. ([BTC](../../pkg/super_node/btc/http_streamer.go), [ETH](../../pkg/super_node/eth/streamer.go)).
* [Payload Converter](../../pkg/super_node/shared/interfaces.go#L34): Converters raw chain data to an intermediary form prepared for IPFS publishing. ([BTC](../../pkg/super_node/btc/converter.go), [ETH](../../pkg/super_node/eth/converter.go)).
* [IPLD Publisher](../../pkg/super_node/shared/interfaces.go#L39): Publishes the converted data to IPFS, returning their CIDs and associated metadata for indexing. ([BTC](../../pkg/super_node/btc/publisher.go), [ETH](../../pkg/super_node/eth/publisher.go)).
* [CID Indexer](../../pkg/super_node/shared/interfaces.go#L44): Indexes CIDs in Postgres with their associated metadata. This metadata is chain specific and selected based on utility. ([BTC](../../pkg/super_node/btc/indexer.go), [ETH](../../pkg/super_node/eth/indexer.go)).
* [CID Retriever](../../pkg/super_node/shared/interfaces.go#L54): Retrieves CIDs from Postgres by searching against their associated metadata, is used to lookup data to serve API requests/subscriptions. ([BTC](../../pkg/super_node/btc/retriever.go), [ETH](../../pkg/super_node/eth/retriever.go)).
* [IPLD Fetcher](../../pkg/super_node/shared/interfaces.go#L62): Fetches the IPLDs needed to service API requests/subscriptions from IPFS using retrieved CIDS; can route through a IPFS block-exchange to search for objects that are not directly available. ([BTC](../../pkg/super_node/btc/ipld_fetcher.go), [ETH](../../pkg/super_node/eth/ipld_fetcher.go))
* [Response Filterer](../../pkg/super_node/shared/interfaces.go#L49): Filters converted data payloads served to API subscriptions; filters according to the subscriber provided parameters. ([BTC](../../pkg/super_node/btc/filterer.go), [ETH](../../pkg/super_node/eth/filterer.go)).
* [API](https://github.com/ethereum/go-ethereum/blob/master/rpc/types.go#L31): Expose RPC methods for clients to interface with the data. Chain-specific APIs should aim to recapitulate as much of the native API as possible. ([VDB](../../pkg/super_node/api.go), [ETH](../../pkg/super_node/eth/api.go)).


Appropriating the service for a new chain is done by creating underlying types to satisfy these interfaces for
the specifics of that chain.

The service uses these interfaces to operate in any combination of three modes: sync, serve, and backfill.
* Sync: Streams raw chain data at the head, converts and publishes it to IPFS, and indexes the resulting set of CIDs in Postgres with useful metadata.
* BackFill: Automatically searches for and detects gaps in the DB; fetches, converts, publishes, and indexes the data to fill these gaps.
* Serve: Opens up IPC, HTTP, and WebSocket servers on top of the superNode DB and any concurrent sync and/or backfill processes.


These three modes are all operated through a single vulcanizeDB command: `superNode`

## Command

Usage: `./vulcanizedb superNode --config={config.toml}`

Configuration can also be done through CLI options and/or environmental variables.
CLI options can be found using `./vulcanizedb superNode --help`.

## Config

Below is the set of universal config parameters for the superNode command, in .toml form, with the respective environmental variables commented to the side.
This set of parameters needs to be set no matter the chain type.

```toml
[database]
    name     = "vulcanize_public" # $DATABASE_NAME
    hostname = "localhost" # $DATABASE_HOSTNAME
    port     = 5432 # $DATABASE_PORT
    user     = "vdbm" # $DATABASE_USER
    password = "" # $DATABASE_PASSWORD

[ipfs]
    path = "~/.ipfs" # $IPFS_PATH

[superNode]
    chain = "bitcoin" # $SUPERNODE_CHAIN
    server = true # $SUPERNODE_SERVER
    ipcPath = "~/.vulcanize/vulcanize.ipc" # $SUPERNODE_IPC_PATH
    wsPath = "127.0.0.1:8082" # $SUPERNODE_WS_PATH
    httpPath = "127.0.0.1:8083" # $SUPERNODE_HTTP_PATH
    sync = true # $SUPERNODE_SYNC
    workers = 1 # $SUPERNODE_WORKERS
    backFill = true # $SUPERNODE_BACKFILL
    frequency = 45 # $SUPERNODE_FREQUENCY
    batchSize = 1 # $SUPERNODE_BATCH_SIZE
    batchNumber = 50 # $SUPERNODE_BATCH_NUMBER
    validationLevel = 1 # $SUPERNODE_VALIDATION_LEVEL
```

Additional parameters need to be set depending on the specific chain.

For Bitcoin:

```toml
[bitcoin]
    wsPath  = "127.0.0.1:8332" # $BTC_WS_PATH
    httpPath = "127.0.0.1:8332" # $BTC_HTTP_PATH
    pass = "password" # $BTC_NODE_PASSWORD
    user = "username" # $BTC_NODE_USER
    nodeID = "ocd0" # $BTC_NODE_ID
    clientName = "Omnicore" # $BTC_CLIENT_NAME
    genesisBlock = "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f" # $BTC_GENESIS_BLOCK
    networkID = "0xD9B4BEF9" # $BTC_NETWORK_ID
```

For Ethereum:

```toml
[ethereum]
    wsPath  = "127.0.0.1:8546" # $ETH_WS_PATH
    httpPath = "127.0.0.1:8545" # $ETH_HTTP_PATH
```

## Database

Currently, the super node persists all data to a single Postgres database. The migrations for this DB can be found [here](../../db/migrations).
Chain-specific data is populated under a chain-specific schema (e.g. `eth` and `btc`) while shared data- such as the IPFS blocks table- is populated under the `public` schema.
Subsequent watchers which act on the raw chain data should build and populate their own schemas or separate databases entirely.

In the future, we will be moving to a foreign table based architecture wherein a single db is used for shared data while each watcher uses
its own database and accesses and acts on the shared data through foreign tables. Isolating watchers to their own databases will prevent complications and
conflicts between watcher db migrations.


## APIs

The super node provides mutliple types of APIs by which to interface with its data.
More detailed information on the APIs can be found [here](apis.md).

## Resync

A separate command `resync` is available for directing the resyncing of data within specified ranges.
This is useful if we want to re-validate a range of data using a new source or clean out bad/deprecated data.
More detailed information on this command can be found [here](resync.md).

## IPFS Considerations

Currently, the IPLD Publisher and Fetcher use internalized IPFS processes which interface directly with a local IPFS repository.
This circumvents the need to run a full IPFS daemon with a [go-ipld-eth](https://github.com/ipfs/go-ipld-eth) plugin, but can lead to issues
with lock-contention on the IPFS repo if another IPFS process is configured and running at the same $IPFS_PATH. This also necessitates the need for
a locally configured IPFS repository.

Once go-ipld-eth has been updated to work with a modern version of PG-IPFS, an additional option will be provided to direct
all publishing and fetching of IPLD objects through a remote IPFS daemon.