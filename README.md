# Vulcanize DB

[![Join the chat at https://gitter.im/vulcanizeio/VulcanizeDB](https://badges.gitter.im/vulcanizeio/VulcanizeDB.svg)](https://gitter.im/vulcanizeio/VulcanizeDB?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.com/vulcanize/maker-vulcanizedb.svg?token=MKcE2K7CRvKtdxSSnbap&branch=staging)](https://travis-ci.com/vulcanize/maker-vulcanizedb)

## About

Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.

## Dependencies
 - Go 1.11+
 - Postgres 10.6
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8.23+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

## Project Setup

Using Vulcanize for the first time requires several steps be done in order to allow use of the software. The following instructions will offer a guide through the steps of the process:

1. Fetching the project
2. Installing dependencies
3. Configuring shell environment
4. Database setup
5. Configuring synced Ethereum node integration
6. Data syncing

### Installation

In order to fetch the project codebase for local use or modification, install it to your `GOPATH` via:

`go get github.com/vulcanize/vulcanizedb`

Once fetched, dependencies can be installed via `go get` or (the preferred method) at specific versions via `golang/dep`, the prototype golang pakcage manager. Installation instructions are [here](https://golang.github.io/dep/docs/installation.html).

In order to install packages with `dep`, ensure you are in the project directory now within your `GOPATH` (default location is `~/go/src/github.com/vulcanize/vulcanizedb/`) and run:

`dep ensure`

After `dep` finishes, dependencies should be installed within your `GOPATH` at the versions specified in `Gopkg.toml`.

Because we are working with a modified version of the go-ethereum accounts/abi package, after running `dep ensure` you will need to run `git checkout vendor/github/ethereum/go-ethereum/accounts/abi` to checkout the modified dependency.
This is explained in greater detail [here](https://github.com/vulcanize/maker-vulcanizedb/issues/41).

Lastly, ensure that `GOPATH` is defined in your shell. If necessary, `GOPATH` can be set in `~/.bashrc` or `~/.bash_profile`, depending upon your system. It can be additionally helpful to add `$GOPATH/bin` to your shell's `$PATH`.

### Setting up the Database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_public`
1. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1.  Run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    - To rollback a single step: `make rollback NAME=vulcanize_public`
    - To rollback to a certain migration: `make rollback_to MIGRATION=n NAME=vulcanize_public`
    - To see status of migrations: `make migration_status NAME=vulcanize_public`

    * See below for configuring additional environments

### Create a migration file
1. `make new_migration NAME=add_columnA_to_table1`
    - This will create a new timestamped migration file in `db/migrations`
1. Write the migration code in the created file, under the respective `goose` pragma
    - Goose automatically runs each migration in a transaction; don't add `BEGIN` and `COMMIT` statements.

### Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` and `levelDbPath`.
  - `ipcPath` should match the local node's IPC filepath:
      - For Geth:
        - The IPC file is called `geth.ipc`.
        - The geth IPC file path is printed to the console when you start geth.
        - The default location is:
          - Mac: `<full home path>/Library/Ethereum`
          - Linux: `<full home path>/ethereum/geth.ipc`

      - For Parity:
        - The IPC file is called `jsonrpc.ipc`.
        - The default location is:
          - Mac: `<full home path>/Library/Application\ Support/io.parity.ethereum/`
          - Linux: `<full home path>/local/share/io.parity.ethereum/`

  - `levelDbPath` should match Geth's chaindata directory path.
      - The geth LevelDB chaindata path is printed to the console when you start geth.
      - The default location is:
          - Mac: `<full home path>/Library/Ethereum/geth/chaindata`
          - Linux: `<full home path>/ethereum/geth/chaindata`
      - `levelDbPath` is irrelevant (and `coldImport` is currently unavailable) if only running parity.

- See `environments/infura.toml` to configure commands to run against infura, if a local node is unavailable.
- Copy `environments/local.toml.example` to `environments/local.toml` to configure commands to run against a local node such as [Ganache](https://truffleframework.com/ganache) or [ganache-cli](https://github.com/trufflesuite/ganache-clihttps://github.com/trufflesuite/ganache-cli).

### Start syncing with postgres
Syncs VulcanizeDB with the configured Ethereum node, populating blocks, transactions, receipts, and logs.
This command is useful when you want to maintain a broad cache of what's happening on the blockchain.
1. Start Ethereum node (**if fast syncing your Ethereum node, wait for initial sync to finish**)
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

### Alternatively, sync from Geth's underlying LevelDB
Sync VulcanizeDB from the LevelDB underlying a Geth node.
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

### Alternatively, sync in "light" mode
Syncs VulcanizeDB with the configured Ethereum node, populating only block headers.
This command is useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values).
1. Start Ethereum node
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb lightSync --config <config.toml> --starting-block-number <block-number>`

## Start full environment in docker by single command

### Geth Rinkeby

make command        | description
------------------- | ----------------
rinkeby_env_up      | start geth, postgres and rolling migrations, after migrations done starting vulcanizedb container
rinkeby_env_deploy  | build and run vulcanizedb container in rinkeby environment
rinkeby_env_migrate | build and run rinkeby env migrations
rinkeby_env_down    | stop and remove all rinkeby env containers

Success run of the VulcanizeDB container require full geth state sync,
attach to geth console and check sync state:

```bash
$ docker exec -it rinkeby_vulcanizedb_geth geth --rinkeby attach
...
> eth.syncing
false
```

If you have full rinkeby chaindata you can move it to `rinkeby_vulcanizedb_geth_data` docker volume to skip long wait of sync.

## Running the Tests
- `createdb vulcanize_private` will create the test db
- `make migrate NAME=vulcanize_private` will run the db migrations
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run the just the integration tests

## Deploying
1. you will need to make sure you have ssh agent running and your ssh key added to it. instructions [here](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/#your-key-must-be-available-to-ssh-agent)
1. `go get -u github.com/pressly/sup/cmd/sup`
1. `sup staging deploy`

## Contract Watchers
Contract watchers work with a light or full sync vDB to fetch raw ethereum data and execute a set of transformations over them, persisting the output.    

A watcher is composed of at least a fetcher and a transformer or set of transformers, where a fetcher is an interface for retrieving raw Ethereum data from some source (e.g. eth_jsonrpc, IPFS)
and a transformer is an interface for filtering through that raw Ethereum data to extract, process, and persist data for specific contracts or accounts. 

### contractWatcher
The `contractWatcher` command is a built-in generic contract watcher. It can watch any and all events for a given contract provided the contract's ABI is available.
It also provides some state variable coverage by automating polling of public methods, with some restrictions:
1. The method must have 2 or less arguments
2. The method's arguments must all be of type address or bytes32 (hash)
3. The method must return a single value

This command operates in two modes- `light` and `full`- which require a light or full-synced vulcanizeDB, respectively.

This command requires the contract ABI be available on Etherscan if it is not provided in the config file by the user.

If method polling is turned on we require an archival node at the ETH ipc endpoint in our config, whether or not we are operating in `light` or `full` mode.
Otherwise, when operating in `light` mode, we only need to connect to a full node to fetch event logs.

This command takes a config of the form:

```toml
  [database]
    name     = "vulcanize_public"
    hostname = "localhost"
    port     = 5432

  [client]
    ipcPath  = "https://mainnet.infura.io/J5Vd2fRtGsw0zZ0Ov3BL"

  [contract]
    network  = ""
    addresses  = [
        "contractAddress1",
        "contractAddress2"
    ]
    [contract.contractAddress1]
        abi    = 'ABI for contract 1'
        startingBlock = 982463
    [contract.contractAddress2]
        abi    = 'ABI for contract 2'
        events = [
            "event1",
            "event2"
        ]
		eventArgs = [
			"arg1",
			"arg2"
		]
        methods = [
            "method1",
			"method2"
        ]
		methodArgs = [
			"arg1",
			"arg2"
		]
        startingBlock = 4448566
        piping = true
````

- The `contract` section defines which contracts we want to watch and with which conditions.
- `network` is only necessary if the ABIs are not provided and wish to be fetched from Etherscan.
    - Empty or nil string indicates mainnet
    - "ropsten", "kovan", and "rinkeby" indicate their respective networks
- `addresses` lists the contract addresses we are watching and is used to load their individual configuration parameters
- `contract.<contractAddress>` are the sub-mappings which contain the parameters specific to each contract address
    - `abi` is the ABI for the contract; if none is provided the application will attempt to fetch one from Etherscan using the provided address and network
    - `events` is the list of events to watch
        - If this field is omitted or no events are provided then by defualt ALL events extracted from the ABI will be watched
        - If event names are provided then only those events will be watched
    - `eventArgs` is the list of arguments to filter events with
        - If this field is omitted or no eventArgs are provided then by default watched events are not filtered by their argument values
        - If eventArgs are provided then only those events which emit at least one of these values as an argument are watched
    - `methods` is the list of methods to poll
        - If this is omitted or no methods are provided then by default NO methods are polled
        - If method names are provided then those methods will be polled, provided
            1) Method has two or less arguments
            1) Arguments are all of address or hash types
            1) Method returns a single value
    - `methodArgs` is the list of arguments to limit polling methods to
        - If this field is omitted or no methodArgs are provided then by default methods will be polled with every combination of the appropriately typed values that have been collected from watched events
        - If methodArgs are provided then only those values will be used to poll methods
    - `startingBlock` is the block we want to begin watching the contract, usually the deployment block of that contract
    - `piping` is a boolean flag which indicates whether or not we want to pipe return method values forward as arguments to subsequent method calls

At the very minimum, for each contract address an ABI and a starting block number need to be provided (or just the starting block if the ABI can be reliably fetched from Etherscan).
With just this information we will be able to watch all events at the contract, but with no additional filters and no method polling.

#### contractWatcher output

Transformed events and polled method results are committed to Postgres in schemas and tables generated according to the contract abi.      

Schemas are created for each contract using the naming convention `<sync-type>_<lowercase contract-address>`
Under this schema, tables are generated for watched events as `<lowercase event name>_event` and for polled methods as `<lowercase method name>_method`  
The 'method' and 'event' identifiers are tacked onto the end of the table names to prevent collisions between methods and events of the same lowercase name    

Example:

Running `./vulcanizedb contractWatcher --config=./environments/example.toml --mode=light`

Runs our contract watcher in light mode, configured to watch the contracts specified in the config file. Note that
by default we operate in `light` mode but the flag is included here to demonstrate its use.

The example config we link to in this example watches two contracts, the ENS Registry (0x314159265dD8dbb310642f98f50C066173C1259b) and TrueUSD (0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E).

Because the ENS Registry is configured with only an ABI and a starting block, we will watch all events for this contract and poll none of its methods. Note that the ENS Registry is an example
of a contract which does not have its ABI available over Etherscan and must have it included in the config file.

The TrueUSD contract is configured with two events (`Transfer` and `Mint`) and a single method (`balanceOf`), as such it will watch these two events and use any addresses it collects emitted from them
to poll the `balanceOf` method with those addresses at every block. Note that we do not provide an ABI for TrueUSD as its ABI can be fetched from Etherscan.

For the ENS contract, it produces and populates a schema with four tables"
`light_0x314159265dd8dbb310642f98f50c066173c1259b.newowner_event`
`light_0x314159265dd8dbb310642f98f50c066173c1259b.newresolver_event`
`light_0x314159265dd8dbb310642f98f50c066173c1259b.newttl_event`
`light_0x314159265dd8dbb310642f98f50c066173c1259b.transfer_event`

For the TrusUSD contract, it produces and populates a schema with three tables:

`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event`
`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.mint_event`
`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method`

Column ids and types for these tables are generated based on the event and method argument names and types and method return types, resulting in tables such as:

Table "light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event"

|  Column    |         Type          | Collation | Nullable |                                           Default                                           | Storage  | Stats target | Description  
|:----------:|:---------------------:|:---------:|:--------:|:-------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event_id_seq'::regclass) | plain    |              |             |
| header_id  | integer               |           | not null |                                                                                             | plain    |              |             |
| token_name | character varying(66) |           | not null |                                                                                             | extended |              |             |
| raw_log    | jsonb                 |           |          |                                                                                             | extended |              |             |
| log_idx    | integer               |           | not null |                                                                                             | plain    |              |             |
| tx_idx     | integer               |           | not null |                                                                                             | plain    |              |             |
| from_      | character varying(66) |           | not null |                                                                                             | extended |              |             |
| to_        | character varying(66) |           | not null |                                                                                             | extended |              |             |
| value_     | numeric               |           | not null |                                                                                             | main     |              |             |


Table "light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method"

|   Column   |         Type          | Collation | Nullable |                                            Default                                            | Storage  | Stats target | Description |
|:----------:|:---------------------:|:---------:|:--------:|:-------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method_id_seq'::regclass) | plain    |              |             |
| token_name | character varying(66) |           | not null |                                                                                               | extended |              |             |
| block      | integer               |           | not null |                                                                                               | plain    |              |             |
| who_       | character varying(66) |           | not null |                                                                                               | extended |              |             |
| returned   | numeric               |           | not null |                                                                                               | main     |              |             |
  
The addition of '_' after table names is to prevent collisions with reserved Postgres words.

Also notice that the contract address used for the schema name has been down-cased.

### composeAndExecute
The `composeAndExecute` command is used to compose and execute over an arbitrary set of custom transformers.
This is accomplished by generating a Go pluggin which allows our `vulcanizedb` binary to link to external transformers, so
long as they abide by our standard [interfaces](https://github.com/vulcanize/maker-vulcanizedb/tree/compose_and_execute/libraries/shared/transformer).   

This command requires Go 1.11+ and [Go plugins](https://golang.org/pkg/plugin/) only work on Unix based systems.

#### Writing custom transformers
Storage Transformers
   * [Guide](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/factories/storage/README.md)   
   * [Example](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/factories/storage/EXAMPLE.md)   
    
Event Transformers   
   * [Guide](https://github.com/vulcanize/maker-vulcanizedb/blob/event_docs/libraries/shared/factories/README.md)
   * [Example](https://github.com/vulcanize/ens_transformers/tree/working)
   
#### composeAndExecute configuration
A .toml config file is specified when executing the command:
`./vulcanizedb composeAndExecute --config=./environments/config_name.toml`

The config provides information for composing a set of transformers:

```toml
[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    user     = "vulcanize"
    password = "vulcanize"
    port     = 5432

[client]
    ipcPath  = "http://kovan0.vulcanize.io:8545"

[exporter]
    home     = "github.com/vulcanize/vulcanizedb"
    name     = "exampleTransformerExporter"
    save     = false
    transformerNames = [
        "transformer1",
        "transformer2",
        "transformer3",
        "transformer4",
    ]
    [exporter.transformer1]
        path = "path/to/transformer1"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer2]
        path = "path/to/transformer2"
        type = "eth_generic"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer3]
        path = "path/to/transformer3"
        type = "eth_event"
        repository = "github.com/account/repo"
        migrations = "db/migrations"
        rank = "0"
    [exporter.transformer4]
        path = "path/to/transformer4"
        type = "eth_storage"
        repository = "github.com/account2/repo2"
        migrations = "to/db/migrations"
        rank = "1"
```
- `home` is the name of the package you are building the plugin for, in most cases this is github.com/vulcanize/vulcanizedb
- `name` is the name used for the plugin files (.so and .go)   
- `save` indicates whether or not the user wants to save the .go file instead of removing it after .so compilation. Sometimes useful for debugging/trouble-shooting purposes.
- `transformerNames` is the list of the names of the transformers we are composing together, so we know how to access their submaps in the exporter map
- `exporter.<transformerName>`s are the sub-mappings containing config info for the transformers
    - `repository` is the path for the repository which contains the transformer and its `TransformerInitializer`
    - `path` is the relative path from `repository` to the transformer's `TransformerInitializer` directory (initializer package).
        - Transformer repositories need to be cloned into the user's $GOPATH (`go get`)
    - `type` is the type of the transformer; indicating which type of watcher it works with (for now, there are only two options: `eth_event` and `eth_storage`)
        - `eth_storage` indicates the transformer works with the [storage watcher](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/watcher/storage_watcher.go)
         that fetches state and storage diffs from an ETH node (instead of, for example, from IPFS)
        - `eth_event` indicates the transformer works with the [event watcher](https://github.com/vulcanize/maker-vulcanizedb/blob/staging/libraries/shared/watcher/event_watcher.go)
         that fetches event logs from an ETH node
        - `eth_contract` indicates the transformer works with the [contract watcher](https://github.com/vulcanize/maker-vulcanizedb/blob/omni_update/libraries/shared/watcher/generic_watcher.go)
        that is made to work with [contract_watcher pkg](https://github.com/vulcanize/maker-vulcanizedb/tree/staging/pkg/omni)
        based transformers which work with either a light or full sync vDB to watch events and poll public methods ([example](https://github.com/vulcanize/ens_transformers/blob/working/transformers/domain_records/transformer.go))
    - `migrations` is the relative path from `repository` to the db migrations directory for the transformer
    - `rank` determines the order that migrations are ran, with lower ranked migrations running first
        - this is to help isolate any potential conflicts between transformer migrations
        - start at "0" 
        - use strings
        - don't leave gaps
        - transformers with identical migrations/migration paths should share the same rank
- Note: If any of the imported transformers need additional config variables those need to be included as well   

This information is used to write and build a Go plugin which exports the configured transformers.
These transformers are loaded onto their specified watchers and executed.

Transformers of different types can be run together in the same command using a single config file or in separate instances using different config files   

The general structure of a plugin .go file, and what we would see built with the above config is shown below

```go
package main

import (
	interface1 "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	transformer1 "github.com/account/repo/path/to/transformer1"
	transformer2 "github.com/account/repo/path/to/transformer2"
	transformer3 "github.com/account/repo/path/to/transformer3"
	transformer4 "github.com/account2/repo2/path/to/transformer4"
)

type exporter string

var Exporter exporter

func (e exporter) Export() []interface1.EventTransformerInitializer, []interface1.StorageTransformerInitializer, []interface1.ContractTransformerInitializer {
	return []interface1.TransformerInitializer{
            transformer1.TransformerInitializer,
            transformer3.TransformerInitializer,
        },     []interface1.StorageTransformerInitializer{
            transformer4.StorageTransformerInitializer,
        },     []interface1.ContractTransformerInitializer{
            transformer2.TransformerInitializer,
        }
}
```

#### Preparing transformer(s) to work as pluggins for composeAndExecute
To plug in an external transformer we need to:

* Create a [package](https://github.com/vulcanize/ens_transformers/blob/working/transformers/registry/new_owner/initializer/initializer.go)
that exports a variable `TransformerInitializer`, `StorageTransformerInitializer`, or `ContractTransformerInitializer` that are of type [TransformerInitializer](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/transformer/event_transformer.go#L33)
or [StorageTransformerInitializer](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/transformer/storage_transformer.go#L31),
or [ContractTransformerInitializer](https://github.com/vulcanize/maker-vulcanizedb/blob/omni_update/libraries/shared/transformer/contract_transformer.go#L31), respectively
* Design the transformers to work in the context of their [event](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/watcher/event_watcher.go#L83),
[storage](https://github.com/vulcanize/maker-vulcanizedb/blob/compose_and_execute/libraries/shared/watcher/storage_watcher.go#L53),
or [contract](https://github.com/vulcanize/maker-vulcanizedb/blob/omni_update/libraries/shared/watcher/contract_watcher.go#L68) watcher execution modes
* Create db migrations to run against vulcanizeDB so that we can store the transformer output
    * Do not `goose fix` the transformer migrations   
    * Specify migration locations for each transformer in the config with the `exporter.transformer.migrations` fields
    * If the base vDB migrations occupy this path as well, they need to be in their `goose fix`ed form
    as they are [here](https://github.com/vulcanize/vulcanizedb/tree/master/db/migrations)

To update a plugin repository with changes to the core vulcanizedb repository, replace the vulcanizedb vendored in the plugin repo (`plugin_repo/vendor/github.com/vulcanize/vulcanizedb`)
with the newly updated version
* The entire vendor lib within the vendored vulcanizedb needs to be deleted (`plugin_repo/vendor/github.com/vulcanize/vulcanizedb/vendor`)
* These complications arise due to this [conflict](https://github.com/golang/go/issues/20481) between `dep` and Go plugins
