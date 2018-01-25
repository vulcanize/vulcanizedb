# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=GKv2Y33qsFnfYgejjvYx&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

## Development Setup

### Dependencies

 - Go 1.9+
 - Postgres 10 
 - Go Ethereum 
    - https://ethereum.github.io/go-ethereum/downloads/ 
    
### Cloning the Repository

1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
    - By default, `go get` does not work for private GitHub repos. This will fix that.
2. `go get github.com/8thlight/vulcanizedb`
3. `cd $GOPATH/src/github.com/8thlight/vulcanizedb`
4. `dep ensure`

### Setting up the Databases

1. Install Postgres
2. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
3. `createdb vulcanize_private`
4. `cd $GOPATH/src/github.com/8thlight/vulcanizedb`
5. `HOST_NAME=localhost NAME=vulcanize_public PORT=5432 make migrate`
    * See below for configuring additional environments

Adding a new migration: `./scripts/create_migration <migration-name>`

### Building 
1. `make build`

### Creating/Using a Private Blockchain

Syncing the public blockchain takes many hours for the initial sync and will download 20+ GB of data.
Here are some instructions for creating a private blockchain that does not depend on having a network connection.

1. Run `./scripts/setup` to create a private blockchain with a new account.
    * This will result in a warning.
2. Run `./scripts/start_private_blockchain`.
3. Run `godo run -- --environment=private` to start listener.

### Connecting to the Public Blockchain

`./scripts/start_blockchain`

### IPC File Paths

The default location for Ethereum is:
 - `$HOME/Library/Ethereum` for Mac
 - `$HOME/.ethereum` for Ubuntu
 - `$GOPATH/src/gihub.com/8thlight/vulcanizedb/test_data_dir/geth.ipc` for private blockchain.

**Note the location of the ipc file is outputted when you connect to geth. It is needed to for configuration**

## Start Vulcanize DB
1. Start geth
2. In a separate terminal start vulcanize_db
    - `vulcanize sync --config <config.toml> --starting-block-number <block-numbe>`

## Watch contract events
1. Start geth
2. In a separate terminal start vulcanize_db
    - `vulcanize sync --config <config.toml> --starting-block-number <block-numbe>`
3. Create event filter 
    - `vulcanize addFilter --config <config.toml> --filter-filepath <filter.json>`
     
### Configuring Additional Environments

You can create configuration files for additional environments.

 * Among other things, it will require the IPC file path
 * See `environments/private.toml` for an example
 * You will need to do this if you want to run a node connecting to the public blockchain

## Running the Tests

### Unit Tests

1. `go test ./pkg/...`

### Integration Test

In order to run the integration tests, you will need to run them against a real blockchain. At the moment the integration tests require [Geth v1.7.2](https://ethereum.github.io/go-ethereum/downloads/) as they depend on the `--dev` mode, which changed in v1.7.3 

1. Run `./scripts/start_private_blockchain` as a separate process.
2. `go test ./...` to run all tests.
