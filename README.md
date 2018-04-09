# Vulcanize DB

[![Build Status](https://travis-ci.org/vulcanize/VulcanizeDB.svg?branch=master)](https://travis-ci.org/vulcanize/VulcanizeDB)

## Dependencies
 - Go 1.9+
 - Postgres 10
 - Ethereum Node
   - [Go Ethereum](https://ethereum.github.io/go-ethereum/downloads/) (1.8+)
   - [Parity 1.8.11+](https://github.com/paritytech/parity/releases)

## Installation
`go get github.com/vulcanize/vulcanizedb`

## Setting up the Database
1. Install Postgres
1. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
1. `createdb vulcanize_public`
1. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
1. Import the schema: `psql vulcanize_public < db/schema.sql`

   or run the migrations: `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    * See below for configuring additional environments

## Configuration
- To use a local Ethereum node, copy `environments/public.toml.example` to
  `environments/public.toml` and update the `ipcPath` to the local node's IPC filepath:
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

- `environments/infura.toml` can be used to run commands with a full
  archive node, rather than a local quick sync node.

## Start syncing with postgres
Syncs VulcanizeDB with the configured Ethereum node.
1. Start node (**if fast syncing wait for initial sync to finish**)
1. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-number>`

## Running the Tests

### Unit Tests
- `go test ./pkg/...`

### Integration Tests
 - `go test ./...` to run all tests.
