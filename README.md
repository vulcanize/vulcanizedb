# Vulcanize DB

[![Build Status](https://travis-ci.com/8thlight/vulcanizedb.svg?token=3psFYN2533rYjhRbvjte&branch=master)](https://travis-ci.com/8thlight/vulcanizedb)

### Dependencies

 - Go 1.9+
 - Postgres 10 
 - Go Ethereum 
    - https://ethereum.github.io/go-ethereum/downloads/ 
 
### Installation 
    go get github.com/vulcanize/vulcanizedb
    
### Setting up the Databases

1. Install Postgres
2. Create a superuser for yourself and make sure `psql --list` works without prompting for a password.
3. `createdb vulcanize_private`
4. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
5. Import the schema: 
    `psql vulcanize_private < db/schema.sql`
   or run the migrations: 
    `make migrate HOST_NAME=localhost NAME=vulcanize_public PORT=5432`
    * See below for configuring additional environments
    
Adding a new migration: `./scripts/create_migration <migration-name>`

## Start syncing with postgres
1. Start geth
2. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-numbe>`
    
   * see `environments` for example config 

## Watch specific contract events
1. Start geth
2. In a separate terminal start vulcanize_db
    - `vulcanizedb sync --config <config.toml> --starting-block-number <block-numbe>`
3. Create event filter 
    - `vulcanizedb addFilter --config <config.toml> --filter-filepath <filter.json>`
    
   * see `filters` for example filter 
     
## Development Setup

### Cloning the Repository (Private repo only)

1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
    - By default, `go get` does not work for private GitHub repos. This will fix that.
2. `go get github.com/vulcanize/vulcanizedb`
3. `cd $GOPATH/src/github.com/vulcanize/vulcanizedb`
4. `dep ensure`

### Creating/Using a test node

Syncing the against the public network takes many hours for the initial sync and will download 20+ GB of data.
Here are some instructions for creating a private test node that does not depend on having a network connection.

1. Run `./scripts/setup` to create a private blockchain with a new account.
    * This will result in a warning.
2. Run `./scripts/start_private_blockchain`.

### IPC File Paths

The default location for Ethereum is:
 - `$HOME/Library/Ethereum` for Mac
 - `$HOME/.ethereum` for Ubuntu
 - `$GOPATH/src/gihub.com/vulcanize/vulcanizedb/test_data_dir/geth.ipc` for private blockchain.

**Note the location of the ipc file is printed to the console when you start geth. It is needed to for configuration**

## Running the Tests

### Unit Tests

1. `go test ./pkg/...`

### Integration Test

In order to run the integration tests, you will need to run them against a real blockchain. At the moment the integration tests require [Geth v1.7.2](https://ethereum.github.io/go-ethereum/downloads/) as they depend on the `--dev` mode, which changed in v1.7.3 

1. Run `./scripts/start_private_blockchain` as a separate process.
2. `go test ./...` to run all tests.
