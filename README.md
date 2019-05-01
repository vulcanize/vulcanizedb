# Vulcanize DB

[![Build Status](https://travis-ci.org/vulcanize/vulcanizedb.svg?branch=master)](https://travis-ci.org/vulcanize/vulcanizedb)
[![Go Report Card](https://goreportcard.com/badge/github.com/vulcanize/vulcanizedb)](https://goreportcard.com/report/github.com/vulcanize/vulcanizedb)

> Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.


## Table of Contents
1. [Background](../staging/README.md#background)
1. [Dependencies](../staging/README.md#dependencies)
1. [Install](../staging/README.md#install)
1. [Usage](../staging/README.md#usage)
1. [Tests](../staging/README.md#tests)
1. [API](../staging/README.md#API)
1. [Contributing](../staging/README.md#contributing)
1. [License](../staging/README.md#license)


## Background
The same data structures and encodings that make Ethereum an effective and trust-less distributed virtual machine
complicate data accessibility and usability for dApp developers. VulcanizeDB improves Ethereum data accessibility by 
providing a suite of tools to ease the extraction and transformation of data into a more useful state.


## Dependencies
 - Go 1.11+
 - Postgres 10.6
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8.23+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)


## Install
1. [Building the project](../staging/README.md#building-the-project)
1. [Setting up the database](../staging/README.md#setting-up-the-database)
1. [Configuring a synced Ethereum node](../staging/README.md#configuring-a-synced-ethereum-node)

### Building the project
Download the codebase to your local `GOPATH` via:

`go get github.com/vulcanize/vulcanizedb`

Move to the project directory and use [golang/dep](https://github.com/golang/dep) to install the dependencies:

`cd $GOPATH/src/github.com/vulcanize/vulcanizedb`

`dep ensure`

Once the dependencies have been successfully installed, build the executable with:

`make build`

If you are running into issues at this stage, ensure that `GOPATH` is defined in your shell.
If necessary, `GOPATH` can be set in `~/.bashrc` or `~/.bash_profile`, depending upon your system.
It can be additionally helpful to add `$GOPATH/bin` to your shell's `$PATH`.

### Setting up the database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_public`
1. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1.  Run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    - There is an optional var `USER=username` if the database user is not the default user `postgres`
    - To rollback a single step: `make rollback NAME=vulcanize_public`
    - To rollback to a certain migration: `make rollback_to MIGRATION=n NAME=vulcanize_public`
    - To see status of migrations: `make migration_status NAME=vulcanize_public`

    * See below for configuring additional environments
    
In some cases (such as recent Ubuntu systems), it may be necessary to overcome failures of password authentication from localhost. To allow access on Ubuntu, set localhost connections via hostname, ipv4, and ipv6 from peer/md5 to trust in: /etc/postgresql/<version>/pg_hba.conf

(It should be noted that trusted auth should only be enabled on systems without sensitive data in them: development and local test databases)

### Configuring a synced Ethereum node
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


## Usage
Usage is broken up into two processes:

### Data syncing
To provide data for transformations, raw Ethereum data must first be synced into vDB.
This is accomplished through the use of the `headerSync`, `sync`, or `coldImport` commands.
These commands are described in detail [here](../staging/documentation/sync.md).

### Data transformation
Contract watchers use the raw data that has been synced into Postgres to filter out and apply transformations to specific data of interest.

There is a built-in `contractWatcher` command which provides generic transformation of most contract data.
The `contractWatcher` command is described further [here](../staging/documentation/contractWatcher.md).

In many cases a custom transformer or set of transformers will need to be written to provide complete or more comprehensive coverage or to optimize other aspects of the output for a specific end-use.
In this case we have provided the `compose`, `execute`, and `composeAndExecute` commands for running custom transformers from external repositories.

Usage of the `compose`, `execute`, and `composeAndExecute` commands is described further [here](../staging/documentation/composeAndExecute.md).

Documentation on how to build custom transformers to work with these commands can be found [here](../staging/documentation/transformers.md).

## Tests
- Replace the empty `ipcPath` in the `environments/infura.toml` with a path to a full node's eth_jsonrpc endpoint (e.g. local geth node ipc path or infura url)
    - Note: integration tests require configuration with an archival node
- `createdb vulcanize_private` will create the test db
- `make migrate NAME=vulcanize_private` will run the db migrations
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run just the integration tests

## API
[Postgraphile](https://www.graphile.org/postgraphile/) is used to expose GraphQL endpoints for our database schemas, this is described in detail [here](../staging/postgraphile/README.md).


## Contributing
Contributions are welcome! For more on this, please see [here](../staging/documentation/contributing.md).

Small note: If editing the Readme, please conform to the [standard-readme specification](https://github.com/RichardLitt/standard-readme).


## License
[AGPL-3.0](../staging/LICENSE) © Vulcanize Inc