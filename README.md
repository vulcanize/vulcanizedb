# Vulcanize DB

[![Join the chat at https://gitter.im/vulcanizeio/VulcanizeDB](https://badges.gitter.im/vulcanizeio/VulcanizeDB.svg)](https://gitter.im/vulcanizeio/VulcanizeDB?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.com/8thlight/maker-vulcanizedb.svg?token=Wi4xzpyShmtvqatRBWkU&branch=staging)](https://travis-ci.com/8thlight/maker-vulcanizedb)

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
`go get gopkg.in/DataDog/dd-trace-go.v1/ddtrace`

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

## Create a migration file (up and down)
1. ./script/create_migrate create_bite_table

## Configuration
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

## Start syncing with postgres
Syncs VulcanizeDB with the configured Ethereum node, populating blocks, transactions, receipts, and logs.
This command is useful when you want to maintain a broad cache of what's happening on the blockchain.
1. Start Ethereum node (**if fast syncing your Ethereum node, wait for initial sync to finish**)
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

## Alternatively, sync from Geth's underlying LevelDB
Sync VulcanizeDB from the LevelDB underlying a Geth node.
1. Assure node is not running, and that it has synced to the desired block height.
1. Start vulcanize_db
   - `./vulcanizedb coldImport --config <config.toml>`
1. Optional flags:
    - `--starting-block-number <block number>`/`-s <block number>`: block number to start syncing from
    - `--ending-block-number <block number>`/`-e <block number>`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

## Alternatively, sync in "light" mode
Syncs VulcanizeDB with the configured Ethereum node, populating only block headers.
This command is useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values).
1. Start Ethereum node
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb lightSync --config <config.toml> --starting-block-number <block-number>`

## Continuously sync Maker event logs from light sync
Continuously syncs Maker event logs from the configured Ethereum node based on the populated block headers.
This includes logs related to auctions, multi-collateral dai, and price feeds.
This command requires that the `lightSync` process is also being run so as to be able to sync in real time.

1. Start Ethereum node (or plan to configure the commands to point to a remote IPC path).
1. In a separate terminal run the lightSync command (see above).
1. In another terminal window run the continuousLogSync command:
  - `./vulcanizedb continuousLogSync --config <config.toml>`
  - An option `--transformers` flag may be passed to the command to specific which transformers to execute, this will default to all transformers if the flag is not passed.
    - `./vulcanizedb continuousLogSync --config environments/private.toml --transformers="priceFeed"`
    - see the `buildTransformerInitializerMap` method in `cmd/continuousLogSync.go` for available transformers

## Backfill Maker event logs from light sync
Backfills Maker event logs from the configured Ethereum node based on the populated block headers.
This includes logs related to auctions, multi-collateral dai, and price feeds.
This command requires that a light sync (see command above) has previously been run.

_Since auction/mcd contracts have not yet been deployed, this command will need to be run a local blockchain at the moment. As such, a new environment file will need to be added. See `environments/local.toml.example`._

1. Start Ethereum node
1. In a separate terminal run the backfill command:
  - `./vulcanizedb backfillMakerLogs --config <config.toml>`
  
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
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run the just the integration tests
- Note: requires Ganache chain setup and seeded with `flip-kick.js` and `frob.js` (in that order)

## Deploying
1. you will need to make sure you have ssh agent running and your ssh key added to it. instructions [here](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/#your-key-must-be-available-to-ssh-agent)
1. `go get -u github.com/pressly/sup/cmd/sup`
1. `sup staging deploy`

## omniWatcher and lightOmniWatcher 
These commands require a pre-synced (full or light) vulcanizeDB (see above sections) 
 
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

