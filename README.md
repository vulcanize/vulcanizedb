# Vulcanize DB

[![Join the chat at https://gitter.im/vulcanizeio/VulcanizeDB](https://badges.gitter.im/vulcanizeio/VulcanizeDB.svg)](https://gitter.im/vulcanizeio/VulcanizeDB?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/vulcanize/vulcanizedb.svg?branch=master)](https://travis-ci.org/vulcanize/vulcanizedb)

## About

Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.

## Dependencies
 - Go 1.9+
 - Postgres 10
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

## Installation
`go get github.com/vulcanize/vulcanizedb`
`go get gopkg.in/DataDog/dd-trace-go.v1/ddtrace`

## Setting up the Database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_public`
1. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1.  Run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`

    * See below for configuring additional environments

## Create a migration file (up and down)
1. ./script/create_migrate create_bite_table

## Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` and `levelDbPath`.
  - `ipcPath` should match the local node's IPC filepath:
      - when using geth:
        - The IPC file is called `geth.ipc`.
        - The geth IPC file path is printed to the console when you start geth.
        - The default location is:
          - Mac: `$HOME/Library/Ethereum`
          - Linux: `$HOME/.ethereum`

      - when using parity:
        - The IPC file is called `jsonrpc.ipc`.
        - The default location is:
          - Mac: `$HOME/Library/Application\ Support/io.parity.ethereum/`
          - Linux: `$HOME/.local/share/io.parity.ethereum/`
          
  - `levelDbPath` should match Geth's chaindata directory path.
      - The geth LevelDB chaindata path is printed to the console when you start geth.
      - The default location is:
          - Mac: `$HOME/Library/Ethereum/geth/chaindata`
          - Linux: `$HOME/.ethereum/geth/chaindata`
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
   - `./vulcanizedb coldImport --config <config.toml> --starting-block-number <block-number> --ending-block-number <block-number>`
1. Optional flags:
    - `--starting-block-number`/`-s`: block number to start syncing from
    - `--ending-block-number`/`-e`: block number to sync to
    - `--all`/`-a`: sync all missing blocks

## Alternatively, sync in "light" mode
Syncs VulcanizeDB with the configured Ethereum node, populating only block headers.
This command is useful when you want a minimal baseline from which to track targeted data on the blockchain (e.g. individual smart contract storage values).
1. Start Ethereum node
1. In a separate terminal start VulcanizeDB:
    - `./vulcanizedb lightSync --config <config.toml> --starting-block-number <block-number>`

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
- `make test`
- Note: requires Ganache chain setup and seeded with `flip-kick.js` and `frob.js` (in that order)

## Deploying
1. you will need to make sure you have ssh agent running and your ssh key added to it. instructions [here](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/#your-key-must-be-available-to-ssh-agent)
1. `go get -u github.com/pressly/sup/cmd/sup`
1. `sup staging deploy`
