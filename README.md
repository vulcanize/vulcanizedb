# Vulcanize DB

[![Join the chat at https://gitter.im/vulcanizeio/VulcanizeDB](https://badges.gitter.im/vulcanizeio/VulcanizeDB.svg)](https://gitter.im/vulcanizeio/VulcanizeDB?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/vulcanize/vulcanizedb.svg?branch=master)](https://travis-ci.org/vulcanize/vulcanizedb)

## About

Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.

## Dependencies
 - Go 1.11+
 - Postgres 10
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8.18+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

## Project Setup

Using Vulcanize for the first time requires several steps be done in order to allow use of the software. The following instructions will offer a guide through the steps of the process:

1. Fetching the project
2. Installing dependencies
3. Configuring shell environment
4. Database setup
5. Configuring synced Ethereum node integration
6. Data syncing

## Installation

In order to fetch the project codebase for local use or modification, install it to your `GOPATH` via:

`go get github.com/vulcanize/vulcanizedb`

Once fetched, dependencies can be installed via `go get` or (the preferred method) at specific versions via `golang/dep`, the prototype golang pakcage manager. Installation instructions are [here](https://golang.github.io/dep/docs/installation.html).

In order to install packages with `dep`, ensure you are in the project directory now within your `GOPATH` (default location is `~/go/src/github.com/vulcanize/vulcanizedb/`) and run:

`dep ensure`

After `dep` finishes, dependencies should be installed within your `GOPATH` at the versions specified in `Gopkg.toml`.

Lastly, ensure that `GOPATH` is defined in your shell. If necessary, `GOPATH` can be set in `~/.bashrc` or `~/.bash_profile`, depending upon your system. It can be additionally helpful to add `$GOPATH/bin` to your shell's `$PATH`.

## Setting up the Database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. Execute `createdb vulcanize_public`
1. Execute `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1. Run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=<postgres port, default 5432>`

    * See below for configuring additional environments

In some cases (such as recent Ubuntu systems), it may be necessary to overcome failures of password authentication from `localhost`. To allow access on Ubuntu, set localhost connections via hostname, ipv4, and ipv6 from `peer`/`md5` to `trust` in: `/etc/postgresql/<version>/pg_hba.conf`

(It should be noted that trusted auth should only be enabled on systems without sensitive data in them: development and local test databases.)

## Configuring Ethereum Node Integration
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

- See `environments/infura.toml` to configure commands to run against infura, if a local node is unavailable. (Support is currently experimental, at this time.)

## Start syncing with postgres
Syncs VulcanizeDB with the configured Ethereum node.
1. Start the node
    - If node state is not yet fully synced, Vulcanize will not be able to operate on the fetched data. You will need to wait for the initial sync to finish.
1. Start the vulcanize_db sync or lightSync
    - Execute `./vulcanizedb sync --config <path to config.toml>`
    - Or `./vulcanizedb lightSync --config <path to config.toml>`
    - Or to sync from a specific block: `./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`
    - Or `./vulcanizedb lightSync --config <config.toml> --starting-block-number <block-number>`

## Alternatively, sync from Geth's underlying LevelDB
Sync VulcanizeDB from the LevelDB underlying a Geth node.
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

## Running the Tests

In order to run the full test suite, a test database must be prepared. By default, the rests use a database named `vulcanize_private`. Create the database in Postgres, and run migrations on the new database in preparation for executing tests:

`make migrate HOST_NAME=localhost NAME=vulcanize_private PORT=<postgres port, default 5432>`

Ginkgo is declared as a `dep` package test execution. Linting and tests can be run together via a provided `make` task:

`make test`

Tests can be run directly via Ginkgo in the project's root directory:

`ginkgo -r`

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

## omniWatcher
This command allows for generic watching of any Ethereum contract provided only the contract's address and additional optional filtering information  
Currently the contract's ABI must be available on etherscan or manually added  
This command will index watched events and polled methods using Postgres tables automatically generated from the ABI  
This command requires a pre-synced (full or light) vulcanizeDB (see above sections)  
 
To watch all events of a contract using a light synced vDB:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address>`  
    
Or if you are using a full synced vDB, change the mode to full:  
    - Execute `./vulcanizedb omniWatcher --mode full --config <path to config.toml> --contract-address <contract address>`  
    
To watch contracts on a network other than mainnet, use the network flag:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --network <ropsten, kovan, or rinkeby>`  
    
To watch events starting at a certain block use the starting block flag:
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --starting-block-number <#>`
    
To watch only specified events use the events flag:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --events <EventName1> --events <EventName2>`  
    
To watch events and poll the specified methods with any addresses and hashes emitted by the watched events utilize the methods flag:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --methods <methodName1> --methods <methodName2>`  
    
To watch specified events and poll the specified method with any addresses and hashes emitted by the watched events:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --events <EventName1> --events <EventName2> --methods <methodName>`  
    
To turn on method piping so that values returned from previous method calls are cached and used as arguments in subsequent method calls:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --piping true --contract-address <contract address> --events <EventName1> --events <EventName2> --methods <methodName>`  
    
To watch all types of events of the contract but only persist the ones that emit one of the filtered-for argument values:  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --event-args <arg1> --event-args <arg2>`  
    
To watch all events of the contract but only poll the specified method with specified argument values (if they are emitted from the watched events):  
    - Execute `./vulcanizedb omniWatcher --config <path to config.toml> --contract-address <contract address> --methods <methodName> --method-args <arg1> --method-args <arg2>`  

### omniWatcher output    

Watched events and methods are transformed and persisted in auto-generated Postgres schemas and tables  
Schemas are created for each contract, with the naming convention `<sync-type>_<lowercase contract-address>`   
e.g. for the TrueUSD contract in lightSync mode we would generate a schema named `light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e`  
Under this schema, tables are generated for watched events as `<lowercase event name>_event` and for polled methods as `<lowercase method name>_method`  
e.g. if we watch Transfer and Mint events of the TrueUSD contract and poll its balanceOf method using the addresses we find emitted from those events we produce a schema with three tables:  

`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.transfer_event`  
`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.mint_event`  
`light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method`  

The 'method' and 'event' identifiers are tacked onto the end of the table names to prevent collisions between methods and events of the same name  

Column id's and types are also auto-generated based on the event and method argument names, resulting in tables such as  

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
     
and   

Table "light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method"  

|   Column   |         Type          | Collation | Nullable |                                            Default                                            | Storage  | Stats target | Description |
|:----------:|:---------------------:|:---------:|:--------:|:-------------------------------------------------------------------------------------------:|:--------:|:------------:|:-----------:|
| id         | integer               |           | not null | nextval('light_0x8dd5fbce2f6a956c3022ba3663759011dd51e73e.balanceof_method_id_seq'::regclass) | plain    |              |             |
| token_name | character varying(66) |           | not null |                                                                                               | extended |              |             |
| block      | integer               |           | not null |                                                                                               | plain    |              |             |
| who_       | character varying(66) |           | not null |                                                                                               | extended |              |             |
| returned   | numeric               |           | not null |                                                                                               | main     |              |             |
    
The addition of '_' after table names is to prevent collisions with reserved Postgres words  