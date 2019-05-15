# Vulcanize DB

[![Build Status](https://travis-ci.org/vulcanize/vulcanizedb.svg?branch=master)](https://travis-ci.org/vulcanize/vulcanizedb)
[![Go Report Card](https://goreportcard.com/badge/github.com/vulcanize/vulcanizedb)](https://goreportcard.com/report/github.com/vulcanize/vulcanizedb)

> Vulcanize DB is a set of tools that make it easier for developers to write application-specific indexes and caches for dapps built on Ethereum.


## Table of Contents
1. [Background](#background)
1. [Install](#install)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)


## Background
The same data structures and encodings that make Ethereum an effective and trust-less distributed virtual machine
complicate data accessibility and usability for dApp developers. VulcanizeDB improves Ethereum data accessibility by
providing a suite of tools to ease the extraction and transformation of data into a more useful state, including
allowing for exposing aggregate data from a suite of smart contracts.

VulanizeDB includes processes that sync, transform and expose data. Syncing involves
querying an Ethereum node and then persisting core data into a Postgres database. Transforming focuses on using previously synced data to
query for and transform log event and storage data for specifically configured smart contract addresses. Exposing data is a matter of getting
data from VulcanizeDB's underlying Postgres database and making it accessible.

![VulcanizeDB Overview Diagram](documentation/diagrams/vdb-overview.png)

## Install

1. [Dependencies](#dependencies)
1. [Building the project](#building-the-project)
1. [Setting up the database](#setting-up-the-database)
1. [Configuring a synced Ethereum node](#configuring-a-synced-ethereum-node)

### Dependencies
 - Go 1.11+
 - Postgres 11.2
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8.23+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

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
    
In some cases (such as recent Ubuntu systems), it may be necessary to overcome failures of password authentication from
localhost. To allow access on Ubuntu, set localhost connections via hostname, ipv4, and ipv6 from peer/md5 to trust in: /etc/postgresql/<version>/pg_hba.conf

(It should be noted that trusted auth should only be enabled on systems without sensitive data in them: development and local test databases)

### Configuring a synced Ethereum node
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` and `levelDbPath`.
  - `ipcPath` should match the local node's IPC filepath:
      - For Geth:
        - The IPC file is called `geth.ipc`.
        - The geth IPC file path is printed to the console when you start geth.
        - The default location is:
          - Mac: `<full home path>/Library/Ethereum/geth.ipc`
          - Linux: `<full home path>/ethereum/geth.ipc`
        - Note: the geth.ipc file may not exist until you've started the geth process

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
As mentioned above, VulcanizeDB's processes can be split into three categories: syncing, transforming and exposing data.

### Data syncing
To provide data for transformations, raw Ethereum data must first be synced into VulcanizeDB.
This is accomplished through the use of the `headerSync`, `fullSync`, or `coldImport` commands.
These commands are described in detail [here](documentation/data-syncing.md).

### Data transformation
Data transformation uses the raw data that has been synced into Postgres to filter out and apply transformations to
specific data of interest. Since there are different types of data that may be useful for observing smart contracts, it
follows that there are different ways to transform this data. We've started by categorizing this into Generic and
Custom transformers:

- Generic Contract Transformer: Generic contract transformation can be done using a built-in command,
`contractWatcher`, which transforms contract events provided the contract's ABI is available. It also
provides some state variable coverage by automating polling of public methods, with some restrictions.
`contractWatcher` is described further [here](documentation/generic-transformer.md).

- Custom Transformers: In many cases custom transformers will need to be written to provide
more comprehensive coverage of contract data. In this case we have provided the `compose`, `execute`, and
`composeAndExecute` commands for running custom transformers from external repositories. Documentation on how to write,
build and run custom transformers as Go plugins can be found
[here](documentation/custom-transformers.md).

### Exposing the data
[Postgraphile](https://www.graphile.org/postgraphile/) is used to expose GraphQL endpoints for our database schemas, this is described in detail [here](documentation/postgraphile.md).


### Tests
- Replace the empty `ipcPath` in the `environments/infura.toml` with a path to a full node's eth_jsonrpc endpoint (e.g. local geth node ipc path or infura url)
    - Note: integration tests require configuration with an archival node
- `createdb vulcanize_private` will create the test db
- `make migrate NAME=vulcanize_private` will run the db migrations
- `make test` will run the unit tests and skip the integration tests
- `make integrationtest` will run just the integration tests


## Contributing
Contributions are welcome!

VulcanizeDB follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/1/4/code-of-conduct).

For more information on contributing, please see [here](documentation/contributing.md).

## License
[AGPL-3.0](LICENSE) © Vulcanize Inc